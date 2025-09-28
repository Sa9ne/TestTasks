package server

import (
	"Medods/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func Start() {
	s := gin.Default()

	s.GET("/Patient", handlers.FindPatient)
	s.GET("/Doctor", handlers.AllDoctors)
	s.GET("/BMR", handlers.BMR)
	s.GET("/BMRHistory", handlers.HistoryBMR)
	s.GET("/BMI", handlers.CalculateBMI)

	err := s.Run(":8080")
	if err != nil {
		log.Fatalf("Сервер не запустился по ошибке: %v", err)
	}
}
