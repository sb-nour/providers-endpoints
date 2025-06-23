package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sb-nour/providers-endpoints/lib"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		fmt.Println("❌ SLACK_WEBHOOK_URL environment variable is not set")
		fmt.Println("Please set your Slack webhook URL:")
		fmt.Println("export SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK")
		return
	}

	fmt.Println("🧪 Testing Slack Integration...")
	fmt.Printf("Webhook URL: %s\n", maskURL(webhookURL))

	// Test 1: Basic message
	fmt.Println("\n📤 Sending basic test message...")
	if err := lib.SendSlackNotification(webhookURL, "🧪 Test message from providers-endpoints!"); err != nil {
		fmt.Printf("❌ Failed to send basic message: %v\n", err)
	} else {
		fmt.Println("✅ Basic message sent successfully!")
	}

	// Test 2: Message with specific channel (if configured)
	if testChannel := os.Getenv("SLACK_TEST_CHANNEL"); testChannel != "" {
		fmt.Printf("\n📤 Sending message to specific channel: %s...\n", testChannel)
		if err := lib.SendSlackNotificationWithChannel(webhookURL, "🧪 Test message with channel override!", fmt.Sprintf("#%s", testChannel)); err != nil {
			fmt.Printf("❌ Failed to send message to channel: %v\n", err)
		} else {
			fmt.Println("✅ Channel-specific message sent successfully!")
		}
	} else {
		fmt.Println("\n💡 To test channel-specific messages, set SLACK_TEST_CHANNEL environment variable")
	}

	// Test 3: Error notification (simulated)
	fmt.Println("\n📤 Sending simulated error notification...")
	lib.SendRegionsFetchErrorNotification("Test Provider", fmt.Errorf("simulated error for testing"))
	fmt.Println("✅ Error notification sent!")

	// Test 4: Show environment variables
	fmt.Println("\n📋 Environment Variables:")
	fmt.Printf("SLACK_WEBHOOK_URL: %s\n", maskURL(os.Getenv("SLACK_WEBHOOK_URL")))
	fmt.Printf("SLACK_ERROR_CHANNEL: %s\n", getEnvOrDefault("SLACK_ERROR_CHANNEL", "not set"))
	fmt.Printf("SLACK_CHANGES_CHANNEL: %s\n", getEnvOrDefault("SLACK_CHANGES_CHANNEL", "not set"))
	fmt.Printf("SLACK_TEST_CHANNEL: %s\n", getEnvOrDefault("SLACK_TEST_CHANNEL", "not set"))

	fmt.Println("\n🎉 Slack integration test completed!")
}

// maskURL masks sensitive parts of URLs for logging
func maskURL(url string) string {
	if len(url) <= 20 {
		return "***"
	}
	return url[:15] + "***" + url[len(url)-10:]
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
