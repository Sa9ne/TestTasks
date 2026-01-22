package server

import (
	"TestTaskService/internal/handlers"
	"log"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	// Используем фреймворк Gin
	s := gin.Default()

	// Для Frontend необходимо добавить CORS, но так как его тут нет, добавлю в комментарий
	// s.Use(cors.Default)

	// Добавляем маршрутизацию
	s.POST("/AddNum", handlers.AddNum)

	// Стартуем по порту 8080
	err := s.Run(":8080")
	if err != nil {
		log.Fatal("Failed to create server")
	}
}
