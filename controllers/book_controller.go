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
	// Extract user from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Fetch user with author
	var user models.User
	if err := database.DB.Preload("Author").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Ensure user has an author profile
	if user.Author.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not have an author profile"})
		return
	}

	// Bind input data
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign author and save book
	book.AuthorID = user.Author.ID
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
