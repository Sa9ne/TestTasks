import aiohttp
import time
from database import SessionLocal
from models.models import Price

# Функция для запроса цены btc
async def fetch_price_btc():
  url = "https://test.deribit.com/api/v2/public/get_index_price?index_name=btc_usd"

  async with aiohttp.ClientSession() as session:
    async with session.get(url) as response:
      return await response.json()
    
# Функция для запроса цены eth
async def fetch_price_eth():
  url = "https://test.deribit.com/api/v2/public/get_index_price?index_name=eth_usd"

  async with aiohttp.ClientSession() as session:
    async with session.get(url) as response:
      return await response.json()

# Функция для сохранения цены
def save_price(ticker: str, price: float):
  db = SessionLocal()
  try:
    price_entry = Price (
      ticker=ticker,
      price=price,
      timestamp=int(time.time())
    )
    db.add(price_entry)
    db.commit()
  finally:
    db.close()