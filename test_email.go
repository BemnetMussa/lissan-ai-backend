// test_email.go - Simple test to demonstrate email functionality
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"lissanai.com/backend/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Create email service
	emailService := service.NewEmailService()

	// Test email sending
	testEmail := "test@example.com"
	testToken := "sample-reset-token-123"
	testName := "John Doe"

	fmt.Println("🧪 Testing email functionality...")
	fmt.Printf("📧 Sending password reset email to: %s\n", testEmail)

	err := emailService.SendPasswordResetEmail(testEmail, testToken, testName)
	if err != nil {
		fmt.Printf("❌ Email sending failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Email functionality test completed!")
	fmt.Println("\n📝 Note:")
	fmt.Println("- If SMTP is configured: Check your email inbox")
	fmt.Println("- If SMTP is NOT configured: Check console output above for reset link")
}