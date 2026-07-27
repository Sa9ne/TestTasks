import asyncio
import io
import os
import zipfile
from contextlib import asynccontextmanager
from datetime import datetime, timezone, timedelta
from pathlib import Path
from typing import Optional, List

import asyncpg
import httpx
from fastapi import FastAPI, HTTPException
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel


API_BASE_URL = os.environ.get("API_BASE_URL", "http://localhost:9000")
CANDIDATE_ID = os.environ.get("CANDIDATE_ID")

DATABASE_URL = os.environ.get(
    "DATABASE_URL", "postgresql://fileservice:fileservice@localhost:5432/file_service"
)

NSK_TZ = timezone(timedelta(hours=7))

CREATE_TABLE_SQL = """
CREATE TABLE IF NOT EXISTS downloaded_files (
    name TEXT PRIMARY KEY,
    downloaded_at TIMESTAMPTZ NOT NULL,
    content TEXT NOT NULL
)
"""


def to_nsk_str(dt: datetime) -> str:
    return dt.astimezone(NSK_TZ).strftime("%Y-%m-%d %H:%M:%S")

pool: Optional[asyncpg.Pool] = None
downloaded_names_set: set = set()


class Progress:

    def __init__(self):
        self.status = "idle"
        self.start_time_nsk: Optional[str] = None
        self.total_names_seen = 0
        self.downloaded_count = 0
        self.error: Optional[str] = None

    def as_dict(self):
        return {
            "status": self.status,
            "start_time_nsk": self.start_time_nsk,
            "total_names_seen": self.total_names_seen,
            "downloaded_count": self.downloaded_count,
            "error": self.error,
        }


progress = Progress()
download_task: Optional[asyncio.Task] = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global pool
    pool = await asyncpg.create_pool(DATABASE_URL, min_size=1, max_size=5)
    async with pool.acquire() as conn:
        await conn.execute(CREATE_TABLE_SQL)
        rows = await conn.fetch("SELECT name FROM downloaded_files")
    downloaded_names_set.update(r["name"] for r in rows)
    progress.downloaded_count = len(downloaded_names_set)
    yield
    await pool.close()


app = FastAPI(title="File download & analysis service", lifespan=lifespan)


# HTTP helper with 429 / 403 handling 
async def api_call(client: httpx.AsyncClient, method: str, path: str, **kwargs) -> httpx.Response:
    url = f"{API_BASE_URL}{path}"
    headers = kwargs.pop("headers", {}) or {}
    if CANDIDATE_ID:
        headers["X-Candidate-Id"] = CANDIDATE_ID

    while True:
        resp = await client.request(method, url, headers=headers, **kwargs)

        if resp.status_code == 429:
            wait = float(resp.headers.get("Retry-After", "5"))
            await asyncio.sleep(wait + 0.5)
            continue

        if resp.status_code == 403:
            wait = resp.headers.get("Retry-After", "1800")
            await asyncio.sleep(float(wait) + 1)
            continue

        resp.raise_for_status()
        return resp


# Background download worker
async def run_download():
    progress.status = "running"
    progress.start_time_nsk = datetime.now(NSK_TZ).strftime("%Y-%m-%d %H:%M:%S")
    progress.error = None

    try:
        async with httpx.AsyncClient(timeout=30) as client:
            while True:
                resp = await api_call(client, "GET", "/api/files/names")
                names = resp.json().get("file_names", [])
                if not names:
                    break

                new_names = [n for n in names if n not in downloaded_names_set]
                progress.total_names_seen += len(new_names)

                for i in range(0, len(new_names), 3):
                    chunk = new_names[i:i + 3]
                    if not chunk:
                        continue

                    dl_resp = await api_call(
                        client, "POST", "/api/files/download", json={"file_names": chunk}
                    )
                    zf = zipfile.ZipFile(io.BytesIO(dl_resp.content))

                    async with pool.acquire() as conn:
                        for fname in zf.namelist():
                            content = zf.read(fname).decode("utf-8", errors="replace")
                            now_utc = datetime.now(timezone.utc)
                            await conn.execute(
                                """
                                INSERT INTO downloaded_files (name, downloaded_at, content)
                                VALUES ($1, $2, $3)
                                ON CONFLICT (name)
                                DO UPDATE SET downloaded_at = EXCLUDED.downloaded_at,
                                              content = EXCLUDED.content
                                """,
                                fname, now_utc, content,
                            )
                            downloaded_names_set.add(fname)
                            progress.downloaded_count += 1

                    await api_call(
                        client, "POST", "/api/files/downloaded", json={"file_names": chunk}
                    )

        progress.status = "done"
    except Exception as e: 
        progress.status = "error"
        progress.error = str(e)


# API routes consumed by the frontend
@app.post("/api/start-download")
async def start_download():
    global download_task
    if progress.status == "running":
        return {"ok": False, "message": "Скачивание уже запущено"}
    download_task = asyncio.create_task(run_download())
    return {"ok": True}


@app.get("/api/progress")
async def get_progress():
    return progress.as_dict()


@app.get("/api/downloaded-files")
async def list_downloaded(page: int = 1, page_size: int = 20, sort: str = "desc"):
    order = "DESC" if sort == "desc" else "ASC"
    offset = (page - 1) * page_size

    async with pool.acquire() as conn:
        total = await conn.fetchval("SELECT COUNT(*) FROM downloaded_files")
        rows = await conn.fetch(
            f"SELECT name, downloaded_at FROM downloaded_files ORDER BY downloaded_at {order} "
            f"LIMIT $1 OFFSET $2",
            page_size, offset,
        )

    items = [{"name": r["name"], "downloaded_at": to_nsk_str(r["downloaded_at"])} for r in rows]
    return {"total": total, "page": page, "page_size": page_size, "items": items}


class CalcRequest(BaseModel):
    names: List[str]


@app.post("/api/calculate")
async def calculate(req: CalcRequest):
    if not req.names:
        raise HTTPException(400, "Список файлов пуст")

    async with pool.acquire() as conn:
        rows = await conn.fetch(
            "SELECT name, content FROM downloaded_files WHERE name = ANY($1::text[])",
            req.names,
        )
    found = {r["name"]: r["content"] for r in rows}
    missing = [n for n in req.names if n not in found]
    if missing:
        raise HTTPException(404, f"Файлы не были скачаны: {', '.join(missing)}")

    total_counts = {str(d): 0 for d in range(10)}
    per_file = {}

    for name in req.names:
        content = found[name].strip()
        counts = {str(d): 0 for d in range(10)}
        for ch in content:
            if ch in counts:
                counts[ch] += 1
                total_counts[ch] += 1
        per_file[name] = counts

    return {"total": total_counts, "per_file": per_file}

app.mount("/", StaticFiles(directory=str(Path(__file__).parent / "static"), html=True), name="static")
