import os
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, declarative_base
from dotenv import load_dotenv

# Загружаем .env файл
load_dotenv()
DATABASE_URL = os.getenv("DATABASE_URL")

# Подключаемся к бд
engine = create_engine(DATABASE_URL)
# Настраиваем orm для работы с бд
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
# Фундамент для всех таблиц
Base = declarative_base()