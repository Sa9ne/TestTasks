package handlers

import (
	"Medods/internal/database"
	"Medods/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HistoryBMR(ctx *gin.Context) {
	id := ctx.Query("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Id не был указан"})
		return
	}

	// пагинация
	limitStr := ctx.Query("limit")
	offsetStr := ctx.Query("offset")

	// Значение по умолчанию
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

	// Ищем в истории BMR с пагинацией
	var BMR []models.BMR
	if err := database.DB.Where("patient_id = ?", id).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&BMR).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Не было найдено пациента"})
		return
	}

	if len(BMR) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"Ошибка": "Пациентов не было найдено"})
		return
	}

	ctx.JSON(http.StatusOK, BMR)
}
