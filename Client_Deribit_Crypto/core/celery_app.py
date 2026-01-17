from celery import Celery

celery_app = Celery(
  "worker",
  broker="redis://localhost:6379/0",
  backend="redis://localhost:6379/0",
)

celery_app.autodiscover_tasks(['core'])

celery_app.conf.beat_schedule = {
  "fetch-btc-eth-every-minute": {
    "task": "core.tasks.fetch_and_save_prices",
    "schedule": 60.0,
  },
}

celery_app.conf.timezone = "UTC"
