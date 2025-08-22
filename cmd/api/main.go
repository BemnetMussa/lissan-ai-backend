// cmd/api/main.go
package main

import (
	"os"
	"lissanai.com/backend/internal/server"
)

// @title           LissanAI Professional API
// @version         1.0
// @description     This is the production-ready foundation for the LissanAI backend.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	server := server.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if the PORT environment variable is not set
	}
	server.Run(":" + port)
}