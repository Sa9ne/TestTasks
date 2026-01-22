package handlers

import (
	"TestTaskService/internal/database"
	"TestTaskService/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddNum(ctx *gin.Context) {
	var Num models.AllNumb

	// Парсим запрос пользователя
	if err := ctx.BindJSON(&Num); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем число в бд
	if err := database.DB.Create(&Num).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении в бд"})
		return
	}

	// Достаем всю информацию о массиве из бд
	var All []models.AllNumb
	if err := database.DB.Order("val asc").Find(&All).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки из бд"})
	}

	// формируем ответ в массиве int
	result := make([]int, len(All))
	for i, n := range All {
		result[i] = n.Val
	}

	// Выводим код протокола 200 и результат
	ctx.JSON(http.StatusOK, result)
}
