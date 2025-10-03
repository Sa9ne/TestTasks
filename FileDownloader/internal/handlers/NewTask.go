package handlers

import (
	"FileDownloader/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Функция для создания уникального ID
func UniqueId() int {
	return int(time.Now().UnixNano())
}

// Функция для сохранения state.json
func Loader(task models.Tasks) {
	// Базовая директория
	wd, _ := os.Getwd()
	baseDir := filepath.Join(wd, "..", "downloads")
	// Путь: downloads/Name_Id
	taskDir := filepath.Join(baseDir, fmt.Sprintf("%s_%d", task.Name, task.Id))

	// Создаем папку
	err := os.MkdirAll(taskDir, os.ModePerm)
	if err != nil {
		log.Printf("не удалось создать папку: %v", err)
		return
	}

	// Путь к state.json
	stateFile := filepath.Join(taskDir, "state.json")

	// Преобразуем task в JSON
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		log.Printf("не удалось создать папку: %v", err)
		return
	}

	// Записываем в файл
	err = os.WriteFile(stateFile, data, 0644)
	if err != nil {
		log.Printf("не удалось создать папку: %v", err)
		return
	}
}

// Скачивание ссылок
func DownloadFiles(task *models.Tasks) {
	// Переходим к директории задачи
	wd, _ := os.Getwd()
	baseDir := filepath.Join(wd, "..", "downloads")
	taskDir := filepath.Join(baseDir, fmt.Sprintf("%s_%d", task.Name, task.Id))

	// Создание wait group для синхронизации выполнения goroutine
	wg := sync.WaitGroup{}
	wg.Add(len(task.Links))

	for i, link := range task.Links {
		// берём расширение из URL
		ext := filepath.Ext(link)
		if ext == "" {
			ext = ".dat" // дефолт, если расширения нет
		}

		filename := fmt.Sprintf("file_%d%s", i+1, ext)
		savePath := filepath.Join(taskDir, filename)

		go func(url, path string) {
			defer wg.Done()

			// Скачиваем содержимое
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Ошибка скачивания %s: %v", url, err)
				return
			}
			defer resp.Body.Close()

			// Проверяем код ответа
			if resp.StatusCode != http.StatusOK {
				log.Printf("Неуспешный ответ %d при скачивании %s", resp.StatusCode, url)
				return
			}

			// Создаём файл для записи
			out, err := os.Create(path)
			if err != nil {
				log.Printf("Ошибка создания файла %s: %v", path, err)
				return
			}
			defer out.Close()

			// Копируем тело ответа в файл
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Printf("Ошибка записи файла %s: %v", path, err)
				return
			}

			log.Printf("Файл сохранён: %s", path)
		}(link, savePath)
	}

	// Ждём окончания всех загрузок
	wg.Wait()

	// Обновляем статус загрузки
	task.Status = "Done"
	Loader(*task)

	// Сообщение пользователю
	log.Printf("Все файлы для задачи %s_%d скачаны", task.Name, task.Id)
}

func NewTask(ctx *gin.Context) {
	// Создаем структуру ввода
	var input struct {
		Name  string   `json:"name"`
		Links []string `json:"links"`
	}

	// Парсим данные в JSON
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверно введенные данные"})
		return
	}

	// Создаем переменную task с введенной информацией пользователя
	task := models.Tasks{
		Id:     UniqueId(),
		Name:   input.Name,
		Links:  input.Links,
		Status: "Processing",
	}

	Loader(task)
	go DownloadFiles(&task)
	ctx.JSON(http.StatusOK, gin.H{"Message": "Task started", "id": task.Id})
}
