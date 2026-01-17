from fastapi import FastAPI
from routes.btc_price import btc_price_router
from routes.eth_price import eth_price_router
from database import Base, engine
from models import Price

# Создаем клиент FastAPI
app = FastAPI()

# создаём таблицы всех моделей автоматически при старте сервера
Base.metadata.create_all(bind=engine)
print("\n Таблицы были успешно созданы! \n")

# Подключаем маршрутизацию
app.include_router(btc_price_router)
app.include_router(eth_price_router)