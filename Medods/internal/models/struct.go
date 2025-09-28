package models

import "time"

// Структура таблицы пациент
type Patient struct {
	Id          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	First_name  string    `json:"first_name" gorm:"unique"`
	Last_name   string    `json:"last_name" gorm:"unique"`
	Middle_name string    `json:"middle_name" gorm:"unique"`
	Birthday    time.Time `json:"birthday" gorm:"type:date;unique"`
	Gender      string    `json:"gender" gorm:"type:char(1)"`
	Height      float64   `json:"height"`
	Weight      float64   `json:"weight"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`

	// Связь таблиц
	Doctors []Doctor `json:"doctors" gorm:"many2many:patient_doctors;"`
	BMRs    []BMR    `gorm:"foreignKey:PatientId"`
}

// Структуру врача
// Убрал уникальность ФИО, так как в задании под пунктом 1.2 не было указана уникальность
type Doctor struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	First_name  string `json:"first_name"`
	Last_name   string `json:"last_name"`
	Middle_name string `json:"middle_name"`

	// Связь таблиц
	Patients []Patient `json:"patients" gorm:"many2many:patient_doctors;"`
}

// Структура BMR
type BMR struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	PatientId int       `json:"patient_id" gorm:"not null"`
	Formula   string    `json:"formula" gorm:"not null"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}
