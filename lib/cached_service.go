package lib

import (
	"fmt"
	"log"

	"github.com/sb-nour/providers-endpoints/service"
)

// CachedProviderFunction wraps a provider function with caching and notification logic
func CachedProviderFunction(providerName string, originalFunc func() service.Regions) func() service.Regions {
	return func() service.Regions {
		// Try to get cached regions first
		cachedRegions, found, err := GetCachedRegions(providerName)
		if err != nil {
			log.Printf("Error checking cache for provider %s: %v", providerName, err)
		}

		// If cache hit and not expired, return cached data
		if found && cachedRegions != nil {
			log.Printf("Using cached regions for provider: %s", providerName)
			return *cachedRegions
		}

		// Cache miss or expired, fetch fresh data
		log.Printf("Cache miss for provider %s, fetching fresh data", providerName)

		var newRegions service.Regions
		var fetchErr error

		// Use a panic recovery mechanism since many providers use panic for errors
		func() {
			defer func() {
				if r := recover(); r != nil {
					fetchErr = fmt.Errorf("provider function panicked: %v", r)
				}
			}()
			newRegions = originalFunc()
		}()

		// Handle fetch errors
		if fetchErr != nil {
			log.Printf("Failed to fetch regions for provider %s: %v", providerName, fetchErr)
			SendRegionsFetchErrorNotification(providerName, fetchErr)

			// Return cached data if available, even if expired
			if cachedRegions != nil {
				log.Printf("Returning expired cached data for provider %s due to fetch error", providerName)
				return *cachedRegions
			}

			// Return empty regions if no cache available
			return service.Regions{
				Storage: make(map[string]string),
				Compute: make(map[string]string),
			}
		}

		// Check if regions have changed (if we have cached data)
		if cachedRegions != nil {
			changed, err := CheckRegionsChanged(providerName, newRegions)
			if err != nil {
				log.Printf("Error checking if regions changed for provider %s: %v", providerName, err)
			} else if changed {
				log.Printf("Regions changed for provider: %s", providerName)
				SendRegionsChangedNotification(providerName, *cachedRegions, newRegions)
			}
		}

		// Cache the new regions
		if err := CacheRegions(providerName, newRegions); err != nil {
			log.Printf("Failed to cache regions for provider %s: %v", providerName, err)
		}

		return newRegions
	}
}

// GetRegionsWithCache is a cached version of GetRegions that uses Turso DB and Slack notifications
func GetRegionsWithCache() map[string]service.Regions {
	// Initialize Turso DB
	if err := InitTursoDB(); err != nil {
		log.Printf("Failed to initialize Turso DB: %v", err)
		log.Printf("Falling back to non-cached mode")
		return GetRegions() // Fall back to original function
	}
	defer func() {
		if err := CloseTursoDB(); err != nil {
			log.Printf("Error closing Turso DB: %v", err)
		}
	}()

	// Create cached versions of provider functions
	cachedProviders := []struct {
		name string
		fn   func() service.Regions
	}{
		{"Amazon AWS", CachedProviderFunction("Amazon AWS", service.GetAmazonRegions)},
		{"Amazon Lightsail", CachedProviderFunction("Amazon Lightsail", service.GetLightsailRegions)},
		{"DigitalOcean", CachedProviderFunction("DigitalOcean", service.GetDigitalOceanRegions)},
		{"UpCloud", CachedProviderFunction("UpCloud", service.GetUpcloudRegions)},
		{"Exoscale", CachedProviderFunction("Exoscale", service.GetExoscaleRegions)},
		{"Google Cloud", CachedProviderFunction("Google Cloud", service.GetGoogleCloudRegions)},
		{"Backblaze", CachedProviderFunction("Backblaze", service.GetBackblazeRegions)},
		{"Linode", CachedProviderFunction("Linode", service.GetLinodeRegions)},
		{"Outscale", CachedProviderFunction("Outscale", service.GetOutscaleRegions)},
		{"Storj", CachedProviderFunction("Storj", service.GetStorjRegions)},
		{"Vultr", CachedProviderFunction("Vultr", service.GetVultrRegions)},
		{"Hetzner", CachedProviderFunction("Hetzner", service.GetHetznerRegions)},
		{"Synology", CachedProviderFunction("Synology", service.GetSynologyRegions)},
	}

	// Use the same concurrent execution pattern as the original GetRegions
	workerCount := 10
	regions := make(map[string]service.Regions)

	// Create channels for communication
	providerRegions := make(chan service.ProviderRegions, len(cachedProviders))
	workerPool := make(chan struct{}, workerCount)

	// Start goroutines for each provider
	for _, provider := range cachedProviders {
		workerPool <- struct{}{}
		go func(provider struct {
			name string
			fn   func() service.Regions
		}) {
			defer func() {
				<-workerPool
			}()
			providerRegions <- service.ProviderRegions{
				Provider: provider.name,
				Regions:  provider.fn(),
			}
		}(provider)
	}

	// Collect results
	for i := 0; i < len(cachedProviders); i++ {
		pr := <-providerRegions
		regions[pr.Provider] = pr.Regions
	}

	close(providerRegions)

	return regions
}

// LogCacheStats logs cache statistics for debugging
func LogCacheStats() {
	if db == nil {
		log.Printf("Database not initialized, cannot show cache stats")
		return
	}

	query := `
		SELECT 
			provider,
			created_at,
			expires_at,
			CASE WHEN expires_at > datetime('now') THEN 'Valid' ELSE 'Expired' END as status
		FROM provider_regions_cache 
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying cache stats: %v", err)
		return
	}
	defer rows.Close()

	log.Printf("=== Cache Statistics ===")
	for rows.Next() {
		var provider, createdAt, expiresAt, status string
		if err := rows.Scan(&provider, &createdAt, &expiresAt, &status); err != nil {
			log.Printf("Error scanning cache stats row: %v", err)
			continue
		}
		log.Printf("Provider: %s | Created: %s | Expires: %s | Status: %s",
			provider, createdAt, expiresAt, status)
	}
	log.Printf("========================")
}
