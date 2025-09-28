package main

import (
	"Medods/internal/database"
	"Medods/internal/server"
)

func main() {
	database.ConnectDB() // Подключение к базе данных
	server.Start()       // Запуск сервера
}
