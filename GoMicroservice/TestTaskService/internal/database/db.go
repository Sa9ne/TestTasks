package database

import (
	"TestTaskService/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Используем библиотеку GORM, как систему ORM, для работы с бд
var DB *gorm.DB

func Connect() {
	// Читаем пути в .env файлах
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Проверяем нашлась ли Database_url
		log.Fatal("Database_url not found")
	}

	// Подключаемся к postgreSQL
	var errOpen error
	DB, errOpen = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errOpen != nil {
		// Обрабатываем ошибку, если подключение не удалось
		log.Fatalf("Failed connect database: %v", errOpen)
	}

	// Проводим автомиграцию таблицы в бд
	DB.AutoMigrate(models.AllNumb{})
}
