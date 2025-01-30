package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title      string `gorm:"not null"`
	AuthorID   uint
	Author     Author
	CategoryID uint
	Category   Category
}
