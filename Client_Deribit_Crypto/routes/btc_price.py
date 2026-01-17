from fastapi import APIRouter
from client.deribit_client import fetch_price_btc

btc_price_router = APIRouter()

@btc_price_router.get("/btc_price")
async def get_btc_price():
  data = await fetch_price_btc()
  return data