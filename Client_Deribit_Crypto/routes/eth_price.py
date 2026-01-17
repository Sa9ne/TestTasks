from fastapi import APIRouter
from client.deribit_client import fetch_price_eth

eth_price_router = APIRouter()

@eth_price_router.get("/eth_price")
async def get_price_eth():
  data = await fetch_price_eth()
  return data