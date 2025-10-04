package handlers

import (
	"FileDownloader/internal/models"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func RecoverUnfinishedTasks() {
	wd, _ := os.Getwd()
	baseDir := filepath.Join(wd, "..", "downloads")

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		log.Printf("Ошибка при чтении директории downloads: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		stateFile := filepath.Join(baseDir, entry.Name(), "state.json")
		data, err := os.ReadFile(stateFile)
		if err != nil {
			log.Printf("Не удалось прочитать %s: %v", stateFile, err)
			continue
		}

		var task models.Tasks
		if err := json.Unmarshal(data, &task); err != nil {
			log.Printf("Ошибка парсинга %s: %v", stateFile, err)
			continue
		}

		// Проверяем, есть ли незавершённые файлы
		if task.Status == "Processing" {
			log.Printf("Найдена незавершённая задача: %s_%d — перезапуск", task.Name, task.Id)
			go DownloadFiles(&task)
		}
	}
}
