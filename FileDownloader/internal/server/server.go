package server

import (
	"FileDownloader/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func Start() {
	s := gin.Default()

	s.POST("/NewTask", handlers.NewTask)

	err := s.Run(":8080")
	if err != nil {
		log.Fatalf("Server was not started:%v", err)
	}
}
