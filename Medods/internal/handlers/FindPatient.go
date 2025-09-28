package handlers

import (
	"Medods/internal/database"
	"Medods/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FindPatient(ctx *gin.Context) {
	var patient []models.Patient
	query := database.DB.Model(&models.Patient{})

	// Фильтр по ФИО
	if fullName := ctx.Query("full_name"); fullName != "" {
		likePattern := "%" + strings.ToLower(fullName) + "%"
		query = query.Where(`
			LOWER(first_name) LIKE ? OR
			LOWER(last_name) LIKE ? OR
			LOWER(middle_name) LIKE ?`, likePattern, likePattern, likePattern)
	}

	// Фильтр по полу
	if gender := ctx.Query("gender"); gender != "" {
		query = query.Where("gender = ?", gender)
	}

	// Фильтр по возрасту
	StartAgeStr := ctx.Query("start_age")
	EndAgeStr := ctx.Query("end_age")

	if StartAgeStr != "" && EndAgeStr != "" {
		// Переводим из строки в число
		StartAge, err1 := strconv.Atoi(StartAgeStr)
		EndAge, err2 := strconv.Atoi(EndAgeStr)

		if err1 == nil && err2 == nil {
			if StartAge > EndAge {
				ctx.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Некорректный диапазон возраста"})
				return
			}
			query = query.Where("EXTRACT(YEAR FROM age(birthday)) BETWEEN ? AND ?", StartAge, EndAge)
		}
	}

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

	if err := query.Find(&patient).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Ошибка": err.Error()})
		return
	}

	if len(patient) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"Ошибка": "Пациенты не найдены"})
		return
	}

	ctx.JSON(http.StatusOK, patient)
}
