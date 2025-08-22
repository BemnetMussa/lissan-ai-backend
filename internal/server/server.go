// internal/server/server.go
package server

import (
	"github.com/gin-gonic/gin"
	"lissanai.com/backend/internal/handler"
	"lissanai.com/backend/internal/repository"
	"lissanai.com/backend/internal/usecase"
	
	_ "lissanai.com/backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	router := gin.Default()

	// --- Dependency Injection ---
	userRepo := repository.NewInMemoryUserRepository()
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	// --- Routes ---
	apiV1 := router.Group("/api/v1")
	{
		auth := apiV1.Group("/auth")
		{
			auth.POST("/signup", authHandler.SignUp)
		}
		// The Interview Team will add their routes here.
	}

	// --- Swagger ---
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	return router
}