// internal/server/server.go
package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	
	"lissanai.com/backend/internal/database"
	"lissanai.com/backend/internal/handler"
	"lissanai.com/backend/internal/middleware"
	"lissanai.com/backend/internal/repository"
	"lissanai.com/backend/internal/service"
	"lissanai.com/backend/internal/usecase"
	
	_ "lissanai.com/backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	router := gin.Default()

	// --- CORS Middleware ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace "*" with your frontend URL in production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- Database Connection ---
	db, err := database.NewMongoConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// --- Services ---
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production" // Default for development
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	
	jwtService := service.NewJWTService(jwtSecret)
	passwordService := service.NewPasswordService()

	// --- Repositories ---
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// --- Use Cases ---
	authUsecase := usecase.NewAuthUsecase(userRepo, refreshTokenRepo, passwordResetRepo, jwtService, passwordService)
	userUsecase := usecase.NewUserUsecase(userRepo, refreshTokenRepo)

	// --- Handlers ---
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	// --- Middleware ---
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// --- Routes ---
	apiV1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/social", authHandler.SocialAuth)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			
			// Protected auth routes
			auth.POST("/logout", authMiddleware, authHandler.Logout)
		}

		// User routes (protected)
		users := apiV1.Group("/users")
		users.Use(authMiddleware)
		{
			users.GET("/me", userHandler.GetProfile)
			users.PATCH("/me", userHandler.UpdateProfile)
			users.DELETE("/me", userHandler.DeleteAccount)
			users.POST("/me/push-token", userHandler.AddPushToken)
		}

		// Future routes for other features
		// interviews := apiV1.Group("/interviews")
		// grammar := apiV1.Group("/grammar")
		// pronunciation := apiV1.Group("/pronunciation")
		// learning := apiV1.Group("/learning")
	}

	// --- Swagger ---
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	return router
}
