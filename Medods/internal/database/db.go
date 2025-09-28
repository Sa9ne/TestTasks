package database

import (
	"Medods/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Объявляем переменную DB для работы с ORM (В моем случае в GORM)
var DB *gorm.DB

func ConnectDB() {
	// Читаем пути в .env файлах
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Введите DATABASE_URL в env файл для определения базы данных")
	}

	// Подключаемся к базе данных через ORM
	var errOpen error
	DB, errOpen = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errOpen != nil {
		log.Fatalf("Ошибка подключение к базе данных: %v", errOpen)
	}

	// Проводим миграцию таблиц в бд
	DB.AutoMigrate(models.Doctor{}, models.Patient{}, models.BMR{})
}
