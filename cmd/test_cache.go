package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sb-nour/providers-endpoints/lib"
	"github.com/sb-nour/providers-endpoints/service"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	// This is a test script to demonstrate caching functionality
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check if Turso DB is configured
	if os.Getenv("TURSO_DATABASE_URL") == "" {
		log.Printf("TURSO_DATABASE_URL not set. Please configure Turso DB to test caching.")
		log.Printf("Example: TURSO_DATABASE_URL=file:test.db go run cmd/test_cache.go")
		return
	}

	// Initialize Turso DB
	if err := lib.InitTursoDB(); err != nil {
		log.Fatalf("Failed to initialize Turso DB: %v", err)
	}
	defer lib.CloseTursoDB()

	// Test caching with a single provider (AWS)
	log.Printf("Testing cache functionality with AWS provider...")

	// Create a cached version of AWS provider
	cachedAWSFunc := lib.CachedProviderFunction("Amazon AWS Test", service.GetAmazonRegions)

	// First call - should fetch fresh data
	log.Printf("First call (should fetch fresh data):")
	regions1 := cachedAWSFunc()
	log.Printf("Storage regions count: %d, Compute regions count: %d",
		len(regions1.Storage), len(regions1.Compute))

	// Second call - should use cached data
	log.Printf("Second call (should use cached data):")
	regions2 := cachedAWSFunc()
	log.Printf("Storage regions count: %d, Compute regions count: %d",
		len(regions2.Storage), len(regions2.Compute))

	// Show cache statistics
	lib.LogCacheStats()

	log.Printf("Cache test completed successfully!")
}
