# Инструкция по тестированию сервиса FileDownloader

## 📌 Запуск проекта
В диретории FileDownloader/cmd запускаем сервис командой:
```bash
go run main.go
```

## ⚙️ Описание портов

Сервер работает на порту 8080
URL: http://localhost:8080

## 🔍 Проверка эндпоинтов 

1. Создание задачи (POST-запрос):
```bash
http://localhost:8080/NewTask
```
В тело запроса передаем например:
```bash
{
    "Name": "FirstTask",
    "Links": ["https://m.media-amazon.com/images/I/41YyiV9n6RL.jpg", "https://i.pinimg.com/474x/e1/07/23/e10723876cbdd44eaad4181d635d64f1.jpg", "https://i.pinimg.com/736x/7e/0a/34/7e0a34a030e0cb470599701d7a5c618e.jpg"]
}
```

В результате этого запроса в директории FileDownloader появиться папка downloader, в которой создастся папка:
```bash
{Name}_{Id}
```

Name - навзание задачи<br>
ID - уникальный id задачи<br>

В этой папке создастся файл state.json, в котором будет храниться вся инфомарция по задаче, а так же в этой папке будет храниться все, что скачалось
<br>
Ответ от сервера: 
```bash
{"Message":"Task started","id":1759578015139881000}
```
2. Проверка статуса задачи (Get-запрос):
```bash
http://localhost:8080/TaskStatus/TaskId

TaskId - заменить на уникальный ID ващей задачи 
```
Ответ от сервера: 
```bash
{
  "id":1759578015139881000,
  "name":"FirstTask",
  "links":["https://m.media-amazon.com/images/I/41YyiV9n6RL.jpg",
          "https://i.pinimg.com/474x/e1/07/23/e10723876cbdd44eaad4181d635d64f1.jpg",
          "https://i.pinimg.com/736x/7e/0a/34/7e0a34a030e0cb470599701d7a5c618e.jpg"],
  "files":[{
          "url":"https://m.media-amazon.com/images/I/41YyiV9n6RL.jpg","filename":"file_1.jpg","status":"Done"},
          {"url":"https://i.pinimg.com/474x/e1/07/23/e10723876cbdd44eaad4181d635d64f1.jpg","filename":"file_2.jpg","status":"Done"},
          {"url":"https://i.pinimg.com/736x/7e/0a/34/7e0a34a030e0cb470599701d7a5c618e.jpg","filename":"file_3.jpg","status":"Done"}],
  "status":"Done"}
```
