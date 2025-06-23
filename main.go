package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sb-nour/providers-endpoints/lib"
	"github.com/sb-nour/providers-endpoints/service"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check for required environment variables
	checkEnvironmentVariables()

	// Use cached version if Turso DB is configured, otherwise fall back to original
	var regions map[string]service.Regions

	if os.Getenv("TURSO_DATABASE_URL") != "" {
		log.Printf("Using cached regions with Turso DB")
		regions = lib.GetRegionsWithCache()

		// Log cache statistics for debugging
		lib.LogCacheStats()
	} else {
		log.Printf("Turso DB not configured, using original non-cached version")
		regions = lib.GetRegions()
	}

	// Marshal and output the results
	regionsJson, err := json.Marshal(regions)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(regionsJson))
}

func checkEnvironmentVariables() {
	// Check for Turso DB configuration
	if tursoURL := os.Getenv("TURSO_DATABASE_URL"); tursoURL != "" {
		log.Printf("Turso DB URL configured: %s", maskSensitiveURL(tursoURL))

		if authToken := os.Getenv("TURSO_AUTH_TOKEN"); authToken != "" {
			log.Printf("Turso auth token configured: %s", maskToken(authToken))
		} else {
			log.Printf("Warning: TURSO_AUTH_TOKEN not set, using database without authentication")
		}
	} else {
		log.Printf("TURSO_DATABASE_URL not set, caching will be disabled")
	}

	// Check for Slack webhook configuration
	if slackURL := os.Getenv("SLACK_WEBHOOK_URL"); slackURL != "" {
		log.Printf("Slack webhook URL configured: %s", maskSensitiveURL(slackURL))
	} else {
		log.Printf("SLACK_WEBHOOK_URL not set, notifications will be disabled")
	}
}

// maskSensitiveURL masks sensitive parts of URLs for logging
func maskSensitiveURL(url string) string {
	if len(url) <= 20 {
		return "***"
	}
	return url[:10] + "***" + url[len(url)-7:]
}

// maskToken masks tokens for logging
func maskToken(token string) string {
	if len(token) <= 10 {
		return "***"
	}
	return token[:4] + "***" + token[len(token)-4:]
}
