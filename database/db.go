package database

import (
	"github.com/ridhotamma/libraryapp/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	dsn := "host=localhost user=ridhotamma password=secret dbname=libraryapp port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB.AutoMigrate(&models.Author{}, &models.Category{}, &models.Book{}, &models.User{})
	return nil
}
