# Client_Deribit_Crypto

Проект на Python для получения цен на криптовалюты (BTC и ETH) с биржи Deribit и сохранения их в базу данных PostgreSQL с использованием Celery для асинхронной периодической обработки.

---

## Design Decisions

- Использован Celery + Redis для периодического получения цен, чтобы не блокировать основной поток.
- Асинхронные функции `fetch_price_btc` и `fetch_price_eth` используют `asyncio` для одновременных запросов к API Deribit.
- SQLAlchemy используется для ORM и удобного взаимодействия с PostgreSQL.
- Все цены сохраняются с UNIX timestamp, что упрощает фильтрацию по времени.
- FastAPI для внешнего API — легкий, быстрый и удобный фреймворк.

---

## Структура проекта

```
Client_Deribit_Crypto/
│
├── client/                 # Клиент для работы с API Deribit
│   ├── __init__.py
│   └── deribit_client.py   # функции fetch_price_btc, fetch_price_eth, save_price
│
├── core/                   # Основная логика Celery
│   ├── __init__.py
│   ├── celery_app.py       # инициализация Celery и Beat schedule
│   └── tasks.py            # задачи Celery
│
├── models/                 # ORM модели базы данных
│   ├── models.py
│   └── __init__.py
├── routes/                 # Роуты для API
│   ├── btc_price.py
│   ├── eth_price.py
│   └── __init__.py
├── main.py                 # точка входа приложения (если нужна)
├── database.py             # подключение к PostgreSQL через SQLAlchemy
├──__init__.py
├── requirements.txt
└── README.md

```

---

## Архитектура проекта

1. **Celery** – асинхронная система задач:
   - **Worker** – выполняет задачи (получает цены и сохраняет их в базу).  
   - **Beat** – планировщик, который каждые 60 секунд ставит задачу `fetch_and_save_prices` в очередь.  
   - **Redis** – брокер сообщений и backend для Celery.

2. **client/deribit_client.py**:
   - `fetch_price_btc()` – асинхронно получает цену BTC с Deribit.  
   - `fetch_price_eth()` – асинхронно получает цену ETH.  
   - `save_price(ticker, price)` – сохраняет цену в базу данных.

3. **core/tasks.py** – определяет задачу Celery `fetch_and_save_prices`, которая:  
   - Получает цены через API.  
   - Сохраняет их в базу.  
   - Возвращает словарь с последними ценами.

4. **core/celery_app.py** – конфигурация Celery:
   - Подключение к брокеру Redis (`redis://localhost:6379/0`).  
   - Автообнаружение задач (`autodiscover_tasks`).  
   - Настройка расписания задач через `beat_schedule`.

5. **database.py** – SQLAlchemy для работы с PostgreSQL.  

---

## Установка и запуск

1. **Клонируем репозиторий и создаем виртуальное окружение:**

```bash
git clone <repo-url>
cd Client_Deribit_Crypto
python3 -m venv dtaskbot-env
source dtaskbot-env/bin/activate
pip install -r requirements.txt
```

2. **Добавляем .env файл в Client_Deribit_Crypto и создаем саму бд**

```.env
DATABASE_URL=postgresql+psycopg2://postgres:password@localhost:5432/DatabaseName
```

password = Пароль от бд

DatabaseName = Название базы данных

postgres = Пользователь бд

5432 = Порт от бд


3. **Настройка переменной окружения для Python path (чтобы Celery видел модули):**
```
export PYTHONPATH="$PWD"
```

4. **Запуск Redis**

```
redis-server
```

5. **Запуск Celery (Worker + Beat вместе, разработка):**

```
celery -A core.celery_app:celery_app worker -B -l info
```

---

## Пример работы задачи
```
[INFO] Task core.tasks.fetch_and_save_prices[...] received
[INFO] Task core.tasks.fetch_and_save_prices[...] succeeded: {'BTC': 95265.21, 'ETH': 3300.88}
```
