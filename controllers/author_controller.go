package controllers

import (
	"net/http"
	"strconv"

	"github.com/ridhotamma/libraryapp/models"

	"github.com/ridhotamma/libraryapp/database"

	"github.com/gin-gonic/gin"
)

func GetAuthors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var authors []models.Author
	if err := database.DB.Limit(limit).Offset(offset).Find(&authors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, authors)
}

// CreateAuthor creates a new author
func CreateAuthor(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, author)
}

// GetAuthor gets an author by ID
func GetAuthor(c *gin.Context) {
	var author models.Author
	id := c.Param("id")
	if err := database.DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}
	c.JSON(http.StatusOK, author)
}

// UpdateAuthor updates an author
func UpdateAuthor(c *gin.Context) {
	var author models.Author
	id := c.Param("id")
	if err := database.DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Save(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, author)
}

// DeleteAuthor deletes an author
func DeleteAuthor(c *gin.Context) {
	var author models.Author
	id := c.Param("id")
	if err := database.DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}
	if err := database.DB.Delete(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Author deleted successfully"})
}
