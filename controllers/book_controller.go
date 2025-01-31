package controllers

import (
	"net/http"
	"strconv"

	"github.com/ridhotamma/libraryapp/database"

	"github.com/ridhotamma/libraryapp/models"

	"github.com/gin-gonic/gin"
)

// GetBooks gets a paginated list of books
func GetBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var books []models.Book
	if err := database.DB.Preload("Author").Preload("Category").Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func CreateBook(c *gin.Context) {
	// 从请求头中获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	var user models.User
	if err := database.DB.Preload("Author").First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将书籍关联到当前用户对应的作者
	book.AuthorID = user.AuthorProfile.ID

	if err := database.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// GetBook gets a book by ID
func GetBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")
	if err := database.DB.Preload("Author").Preload("Category").First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// UpdateBook updates a book
func UpdateBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")
	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book
func DeleteBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")
	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if err := database.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
