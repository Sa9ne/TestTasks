package main

import (
	"TestTaskService/internal/database"
	"TestTaskService/internal/server"
)

func main() {
	database.Connect()
	server.Start()
}
