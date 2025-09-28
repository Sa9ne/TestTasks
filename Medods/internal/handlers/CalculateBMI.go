package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// структура для ответа от внешнего API
type BMIResponse struct {
	BMI    float64 `json:"bmi"`
	Status string  `json:"status"`
}

func CalculateBMI(ctx *gin.Context) {
	// Чтение параметров (вес и рост)
	weightStr := ctx.Query("weight")
	heightStr := ctx.Query("height")

	if weightStr == "" || heightStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Ошибка": "weight и height обязательны"})
		return
	}

	// Формируем адрес внешнего API
	apiURL := fmt.Sprintf("https://bmicalculatorapi.vercel.app/api/bmi/%s/%s", weightStr, heightStr)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "не удалось создать запрос"})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "не удалось выполнить запрос к API"})
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать ответ API"})
		return
	}

	// Парсим JSON в структуру
	var result BMIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка парсинга ответа API"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
