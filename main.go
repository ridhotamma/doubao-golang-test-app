package main

import (
	"github.com/ridhotamma/libraryapp/controllers"
	"github.com/ridhotamma/libraryapp/database"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/libraryapp/models"
)

func main() {
	if err := database.ConnectDB(); err != nil {
		panic("Failed to connect to database")
	}
	database.DB.AutoMigrate(&models.Author{}, &models.Category{}, &models.Book{})

	r := gin.Default()

	// Author routes
	authors := r.Group("/authors")
	{
		authors.POST("", controllers.CreateAuthor)
		authors.GET("/:id", controllers.GetAuthor)
		authors.GET("", controllers.GetAuthors)
		authors.PUT("/:id", controllers.UpdateAuthor)
		authors.DELETE("/:id", controllers.DeleteAuthor)
	}

	// Category routes
	categories := r.Group("/categories")
	{
		categories.POST("", controllers.CreateCategory)
		categories.GET("/:id", controllers.GetCategory)
		categories.GET("", controllers.GetCategories)
		categories.PUT("/:id", controllers.UpdateCategory)
		categories.DELETE("/:id", controllers.DeleteCategory)
	}

	// Book routes
	books := r.Group("/books")
	{
		books.POST("", controllers.CreateBook)
		books.GET("/:id", controllers.GetBook)
		books.GET("", controllers.GetBooks)
		books.PUT("/:id", controllers.UpdateBook)
		books.DELETE("/:id", controllers.DeleteBook)
	}

	users := r.Group("/users")
	{
		users.POST("", controllers.CreateUser)
		users.GET("/:id", controllers.GetUser)
		users.GET("", controllers.GetUsers)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
		users.POST("/login", controllers.LoginUser)
	}

	r.Run(":8080")
}
