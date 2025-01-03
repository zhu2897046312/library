package router

import (
	_ "library/docs" // 导入 swagger docs
	"library/handler"
	"library/middleware"
	"library/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Library Management System API
// @version 1.0
// @description 图书馆管理系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// SetupRouter initializes the router and sets up all routes
func SetupRouter(factory service.Factory) *gin.Engine {
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORSMiddleware())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create handlers
	userHandler := handler.NewUserHandler(factory.GetUserService())
	bookHandler := handler.NewBookHandler(factory.GetBookService())
	borrowHandler := handler.NewBorrowHandler(factory.GetBorrowService())
	reviewHandler := handler.NewReviewHandler(factory.GetReviewService())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)

			auth := users.Use(middleware.AuthMiddleware())
			{
				auth.GET("/profile", userHandler.GetProfile)
				auth.PUT("/profile", userHandler.UpdateProfile)
				auth.PUT("/password", userHandler.ChangePassword)
			}

			admin := auth.Use(middleware.AdminAuthMiddleware())
			{
				admin.GET("", userHandler.ListUsers)
			}
		}

		// Book routes
		books := v1.Group("/books")
		{
			books.GET("", bookHandler.ListBooks)
			books.GET("/:id", bookHandler.GetBook)

			auth := books.Use(middleware.AuthMiddleware())
			{
				admin := auth.Use(middleware.AdminAuthMiddleware())
				{
					admin.POST("", bookHandler.CreateBook)
					admin.PUT("/:id", bookHandler.UpdateBook)
					admin.PUT("/:id/status", bookHandler.UpdateBookStatus)
					admin.PUT("/:id/stock", bookHandler.UpdateBookStock)
				}
			}
		}

		// Borrow routes
		borrows := v1.Group("/borrows")
		{
			auth := borrows.Use(middleware.AuthMiddleware())
			{
				auth.GET("", borrowHandler.ListBorrows)
				auth.GET("/:id", borrowHandler.GetBorrow)
				auth.POST("", borrowHandler.BorrowBook)
				auth.POST("/return", borrowHandler.ReturnBook)
			}
			admin := auth.Use(middleware.AdminAuthMiddleware())
			{
				admin.PUT("/:id", borrowHandler.UpdateBorrow)
			}
		}

		// Review routes
		reviews := v1.Group("/reviews")
		{
			reviews.GET("", reviewHandler.ListReviews)
			reviews.GET("/:id", reviewHandler.GetReview)

			auth := reviews.Use(middleware.AuthMiddleware())
			{
				auth.POST("", reviewHandler.CreateReview)
				auth.PUT("/:id", reviewHandler.UpdateReview)
				auth.DELETE("/:id", reviewHandler.DeleteReview)

				admin := auth.Use(middleware.AdminAuthMiddleware())
				{
					admin.PUT("/:id/status", reviewHandler.UpdateReviewStatus)
				}
			}
		}



	}

	return r
}
