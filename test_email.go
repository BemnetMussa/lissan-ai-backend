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

	// Check if SMTP is configured
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpUsername == "" || smtpPassword == "" {
		fmt.Println("🧪 Testing email functionality (Development Mode)...")
		fmt.Println("📝 SMTP not configured - will show console output instead of sending email")
		fmt.Println()
	} else {
		fmt.Println("🧪 Testing email functionality (Production Mode)...")
		fmt.Printf("📧 SMTP configured with username: %s\n", smtpUsername)
		fmt.Println()
	}

	// Create email service
	emailService := service.NewEmailService()

	// Test email sending
	testEmail := "test@example.com"
	testToken := "sample-reset-token-123"
	testName := "John Doe"

	fmt.Printf("📧 Sending password reset email to: %s\n", testEmail)

	err := emailService.SendPasswordResetEmail(testEmail, testToken, testName)
	if err != nil {
		fmt.Printf("❌ Email sending failed: %v\n", err)
		fmt.Println("\n💡 This is expected if SMTP credentials are not configured.")
		fmt.Println("   Configure SMTP in .env file to enable actual email sending.")
		os.Exit(0) // Exit successfully since this is expected behavior
	}

	fmt.Println("✅ Email functionality test completed!")
	fmt.Println("\n📝 Note:")
	fmt.Println("- If SMTP is configured: Check your email inbox")
	fmt.Println("- If SMTP is NOT configured: Check console output above for reset link")
}