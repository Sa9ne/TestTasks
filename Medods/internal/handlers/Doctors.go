package handlers

import (
	"Medods/internal/database"
	"Medods/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AllDoctors(ctx *gin.Context) {
	var doctors []models.Doctor
	query := database.DB.Model(&models.Doctor{})

	// Пагинация
	limitStr := ctx.Query("limit")
	offsetStr := ctx.Query("offset")

	// Создаем значения по умолчанию, если они не заданы вручную
	limit := 5
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&doctors).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка при получении врачей"})
		return
	}

	if len(doctors) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"Ошибка": "Врачей не было найдено"})
		return
	}

	ctx.JSON(http.StatusOK, doctors)
}
