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
func Loader(task *models.Tasks) error {
	// Базовая директория
	wd, _ := os.Getwd()
	baseDir := filepath.Join(wd, "..", "downloads")
	// Путь: downloads/Name_Id
	taskDir := filepath.Join(baseDir, fmt.Sprintf("%s_%d", task.Name, task.Id))
	stateFile := filepath.Join(taskDir, "state.json")

	// Создаем папку
	err := os.MkdirAll(taskDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("не удалось создать папку: %v", err)
	}

	// Преобразуем в JSON
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %v", err)
	}

	// Чтобы не потерять файл при сбое — пишем во временный
	tmpFile := stateFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("ошибка записи tmp файла: %v", err)
	}

	// Атомарно заменяем старый state.json новым
	if err := os.Rename(tmpFile, stateFile); err != nil {
		return fmt.Errorf("ошибка замены state файла: %v", err)
	}

	return nil
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

	// Инициализация списка файлов
	task.Files = make([]models.Files, len(task.Links))
	for i, link := range task.Links {
		// берём расширение из URL
		ext := filepath.Ext(link)
		if ext == "" {
			ext = ".dat" // дефолт, если расширения нет
		}
		filename := fmt.Sprintf("file_%d%s", i+1, ext)

		task.Files[i] = models.Files{
			Url:      link,
			FileName: filename,
			Status:   "Processing",
		}

		savePath := filepath.Join(taskDir, filename)

		go func(index int, url, path string) {
			defer wg.Done()

			// Скачиваем содержимое
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Ошибка скачивания %s: %v", url, err)
				task.Files[index].Status = "Error"
				Loader(task)
				return
			}
			defer resp.Body.Close()

			// Проверяем код ответа
			if resp.StatusCode != http.StatusOK {
				log.Printf("Неуспешный ответ %d при скачивании %s", resp.StatusCode, url)
				task.Files[index].Status = "Error"
				Loader(task)
				return
			}

			// Создаём файл для записи
			out, err := os.Create(path)
			if err != nil {
				log.Printf("Ошибка создания файла %s: %v", path, err)
				task.Files[index].Status = "Error"
				Loader(task)
				return
			}
			defer out.Close()

			// Копируем тело ответа в файл
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Printf("Ошибка записи файла %s: %v", path, err)
				task.Files[index].Status = "Error"
				Loader(task)
				return
			}

			// Сохраняем изменения для каждого файла
			task.Files[index].Status = "Done"
			Loader(task)
			log.Printf("Файл сохранён: %s", path)
		}(i, link, savePath)
	}

	// Ждём окончания всех загрузок
	wg.Wait()

	// Обновляем статус загрузки
	task.Status = "Done"
	Loader(task)

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

	// Записываем все ссылки в отдельные структуры, для сохранения статуса каждой из них
	files := make([]models.Files, len(input.Links))
	for i, links := range input.Links {
		files[i] = models.Files{
			Url:      links,
			FileName: fmt.Sprintf("file_%d%s", i+1, filepath.Ext(links)),
			Status:   "Processing",
		}
	}

	// Создаем переменную task с введенной информацией пользователя
	task := models.Tasks{
		Id:     UniqueId(),
		Name:   input.Name,
		Links:  input.Links,
		Files:  files,
		Status: "Processing",
	}

	Loader(&task)
	go DownloadFiles(&task)
	ctx.JSON(http.StatusOK, gin.H{"Message": "Task started", "id": task.Id})
}
