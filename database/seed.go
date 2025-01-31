package database

import (
	"log"

	"github.com/ridhotamma/libraryapp/models"
)

// SeedData 初始化种子数据
func SeedData() {
	// 初始化作者数据
	authors := []models.Author{
		{Name: "Author 1"},
		{Name: "Author 2"},
	}
	for _, author := range authors {
		if err := DB.Create(&author).Error; err != nil {
			log.Printf("Failed to create author: %v", err)
		}
	}

	// 初始化分类数据
	categories := []models.Category{
		{Name: "Category 1"},
		{Name: "Category 2"},
	}
	for _, category := range categories {
		if err := DB.Create(&category).Error; err != nil {
			log.Printf("Failed to create category: %v", err)
		}
	}

	// 初始化用户数据
	users := []models.User{
		{Username: "admin", Password: "admin123"},
		{Username: "user", Password: "user123"},
	}
	for _, user := range users {
		if err := user.SetPassword(user.Password); err != nil {
			log.Printf("Failed to hash password for user %s: %v", user.Username, err)
			continue
		}
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
		}
	}
}
