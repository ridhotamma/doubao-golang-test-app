package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ridhotamma/libraryapp/database"
	"github.com/ridhotamma/libraryapp/models"
)

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var users []models.User
	if err := database.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser 创建新用户
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 对用户密码进行哈希处理
	if err := user.SetPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// GetUser 获取单个用户
func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果更新了密码，需要重新哈希处理
	if user.Password != "" {
		if err := user.SetPassword(user.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user from database
	var storedUser models.User
	if err := database.DB.Where("username = ?", inputUser.Username).Preload("Author").First(&storedUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check password
	if !storedUser.CheckPassword(inputUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Create Author if not exists
	if storedUser.Author.ID == 0 {
		newAuthor := models.Author{
			Name:   storedUser.Username,
			UserID: storedUser.ID,
		}
		if err := database.DB.Create(&newAuthor).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
			return
		}
		storedUser.Author = newAuthor
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   storedUser.ID,
		"author_id": storedUser.Author.ID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login successful",
		"token":     tokenString,
		"author_id": storedUser.Author.ID,
	})
}
