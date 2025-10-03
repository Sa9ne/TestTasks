package handlers

import (
	"FileDownloader/internal/models"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func TaskStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	wd, _ := os.Getwd()
	baseDir := filepath.Join(wd, "..", "downloads")
	entries, _ := os.ReadDir(baseDir)
	for _, e := range entries {
		// Проверка ID
		if strings.HasSuffix(e.Name(), "_"+id) {
			// Идем по директории
			stateFile := filepath.Join(baseDir, e.Name(), "state.json")

			// Читаем содержимое state.json
			data, err := os.ReadFile(stateFile)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось прочитать state.json"})
				return
			}

			// Парсим JSON
			var task models.Tasks
			if err := json.Unmarshal(data, &task); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга state.json"})
				return
			}

			// Успешный вывод
			ctx.JSON(http.StatusOK, task)
			return
		}
	}

	// Ошибка поиска
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
}
