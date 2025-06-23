package lib

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sb-nour/providers-endpoints/service"
	_ "github.com/tursodatabase/libsql-client-go/libsql" // Register libsql driver
)

type CacheEntry struct {
	Provider    string    `json:"provider"`
	RegionsHash string    `json:"regions_hash"`
	Regions     string    `json:"regions"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

const (
	CacheDuration = 24 * time.Hour // Cache for 24 hours
)

var db *sql.DB

func InitTursoDB() error {
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if dbURL == "" {
		return fmt.Errorf("TURSO_DATABASE_URL environment variable is required")
	}

	log.Printf("Attempting to connect to Turso DB: %s", dbURL)

	var err error
	if authToken != "" {
		db, err = sql.Open("libsql", dbURL+"?authToken="+authToken)
	} else {
		db, err = sql.Open("libsql", dbURL)
	}

	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTable(); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func createTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS provider_regions_cache (
			provider TEXT PRIMARY KEY,
			regions_hash TEXT NOT NULL,
			regions TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			expires_at DATETIME NOT NULL
		)
	`
	_, err := db.Exec(query)
	return err
}

func hashRegions(regions service.Regions) string {
	regionsJSON, _ := json.Marshal(regions)
	hash := sha256.Sum256(regionsJSON)
	return hex.EncodeToString(hash[:])
}

func GetCachedRegions(provider string) (*service.Regions, bool, error) {
	if db == nil {
		return nil, false, fmt.Errorf("database not initialized")
	}

	query := `
		SELECT regions_hash, regions, expires_at 
		FROM provider_regions_cache 
		WHERE provider = ? AND expires_at > datetime('now')
	`

	var regionsHash, regionsJSON string
	var expiresAt time.Time

	err := db.QueryRow(query, provider).Scan(&regionsHash, &regionsJSON, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil // Cache miss
		}
		return nil, false, fmt.Errorf("failed to query cache: %w", err)
	}

	var regions service.Regions
	if err := json.Unmarshal([]byte(regionsJSON), &regions); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal cached regions: %w", err)
	}

	return &regions, true, nil
}

func CacheRegions(provider string, regions service.Regions) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	regionsJSON, err := json.Marshal(regions)
	if err != nil {
		return fmt.Errorf("failed to marshal regions: %w", err)
	}

	regionsHash := hashRegions(regions)
	now := time.Now()
	expiresAt := now.Add(CacheDuration)

	query := `
		INSERT OR REPLACE INTO provider_regions_cache 
		(provider, regions_hash, regions, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query, provider, regionsHash, string(regionsJSON), now, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to cache regions: %w", err)
	}

	log.Printf("Cached regions for provider: %s", provider)
	return nil
}

func CheckRegionsChanged(provider string, newRegions service.Regions) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("database not initialized")
	}

	query := `SELECT regions_hash FROM provider_regions_cache WHERE provider = ?`

	var oldHash string
	err := db.QueryRow(query, provider).Scan(&oldHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil // No previous cache, consider it changed
		}
		return false, fmt.Errorf("failed to check regions hash: %w", err)
	}

	newHash := hashRegions(newRegions)
	return oldHash != newHash, nil
}

func CloseTursoDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
