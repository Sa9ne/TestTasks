package handlers

import (
	"Medods/internal/database"
	"Medods/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func BMR(ctx *gin.Context) {
	id := ctx.Query("id")
	formula := ctx.Query("formula")

	// Проверка ввода id и formula
	if id == "" || formula == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Отсутствует ID или Формула"})
		return
	}

	// Проверка id пациента
	PatientId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id"})
		return
	}

	// Поиск пациента в базе
	var patient models.Patient
	if err := database.DB.First(&patient, PatientId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Ошибка": "Пациент не найден"})
		return
	}

	// Узнаем возраст для вычислений
	age := time.Now().Year() - patient.Birthday.Year()

	// Считаем формулы
	var result float64
	switch formula {
	// Формула Миффлина
	case "mifflin":
		s := -161
		if patient.Gender == "м" {
			s = 5
		}
		result = 10*patient.Weight + 6.25*patient.Height - 5*float64(age) + float64(s)
	// Формула Харриса
	case "harris":
		if patient.Gender == "м" {
			result = 88.36 + (13.4 * patient.Weight) + (4.8 * patient.Height) - (5.7 * float64(age))
		} else {
			result = 447.6 + (9.2 * patient.Weight) + (3.1 * patient.Height) - (4.3 * float64(age))
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неизвестная формула"})
		return
	}

	// Сохраняем результат
	bmr := models.BMR{
		PatientId: PatientId,
		Formula:   formula,
		Result:    result,
		CreatedAt: time.Now(),
	}

	// Создаем запись
	if err := database.DB.Create(&bmr).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения результата"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"patient_id": PatientId,
		"formula":    formula,
		"result":     result,
	})
}
