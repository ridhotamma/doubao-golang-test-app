package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ridhotamma/libraryapp/controllers"
	"github.com/ridhotamma/libraryapp/database"
	"github.com/ridhotamma/libraryapp/middlewares"
	"github.com/ridhotamma/libraryapp/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// 连接数据库
	if err := database.ConnectDB(); err != nil {
		panic("Failed to connect to database")
	}
	// 自动迁移数据库表
	database.DB.AutoMigrate(&models.Author{}, &models.Category{}, &models.Book{}, &models.User{})

	// 调用种子数据初始化函数
	database.SeedData()

	// 创建默认的 Gin 引擎
	r := gin.Default()

	// 无需身份验证的路由组
	unprotected := r.Group("/api")
	{
		// 用户登录路由，不需要身份验证
		users := unprotected.Group("/users")
		users.POST("/login", controllers.LoginUser)
	}

	// 需要身份验证的路由组
	protected := r.Group("/api")
	// 应用 JWT 身份验证中间件
	protected.Use(middlewares.AuthMiddleware())
	{
		// 作者相关路由
		authors := protected.Group("/authors")
		{
			authors.POST("", controllers.CreateAuthor)
			authors.GET("/:id", controllers.GetAuthor)
			authors.GET("", controllers.GetAuthors)
			authors.PUT("/:id", controllers.UpdateAuthor)
			authors.DELETE("/:id", controllers.DeleteAuthor)
		}

		// 分类相关路由
		categories := protected.Group("/categories")
		{
			categories.POST("", controllers.CreateCategory)
			categories.GET("/:id", controllers.GetCategory)
			categories.GET("", controllers.GetCategories)
			categories.PUT("/:id", controllers.UpdateCategory)
			categories.DELETE("/:id", controllers.DeleteCategory)
		}

		// 书籍相关路由
		books := protected.Group("/books")
		{
			books.POST("", controllers.CreateBook)
			books.GET("/:id", controllers.GetBook)
			books.GET("", controllers.GetBooks)
			books.PUT("/:id", controllers.UpdateBook)
			books.DELETE("/:id", controllers.DeleteBook)
		}

		// 用户相关路由（除了登录）
		users := protected.Group("/users")
		{
			users.POST("", controllers.CreateUser)
			users.GET("/:id", controllers.GetUser)
			users.GET("", controllers.GetUsers)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
		}
	}

	// 启动服务器
	r.Run(":8080")
}
