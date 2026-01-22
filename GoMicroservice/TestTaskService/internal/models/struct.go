package models

type AllNumb struct {
	ID  int `gorm:"primaryKey"` // Указываем primaryKey
	Val int `gorm:"not null"`   // Указываем что число обязательно
}
