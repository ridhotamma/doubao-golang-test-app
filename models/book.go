package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title      string `gorm:"not null"`
	AuthorID   uint
	Author     Author `gorm:"foreignKey:AuthorID"`
	CategoryID uint
	Category   Category `gorm:"foreignKey:CategoryID"`
}
