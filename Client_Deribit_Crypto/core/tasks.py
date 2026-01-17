from celery import shared_task
from client.deribit_client import fetch_price_btc, fetch_price_eth, save_price
import asyncio

@shared_task(bind=True, name="core.tasks.fetch_and_save_prices")
def fetch_and_save_prices(self):
  async def fetch():
    btc = await fetch_price_btc()
    eth = await fetch_price_eth()
    return btc, eth

  btc_data, eth_data = asyncio.run(fetch())

  btc_price = btc_data["result"]["index_price"]
  eth_price = eth_data["result"]["index_price"]
  
  save_price("BTC", btc_price)
  save_price("ETH", eth_price)

  return {
    "BTC": btc_price,
    "ETH": eth_price,
  }
