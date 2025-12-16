package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edvirons/ssp/ims/internal/seed"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Command line flags
	clean := flag.Bool("clean", false, "Clean (delete) existing seed data before seeding")
	tenant := flag.String("tenant", "demo-tenant", "Tenant ID for seed data")
	verbose := flag.Bool("verbose", true, "Enable verbose logging")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without making changes")

	flag.Parse()

	// Get database URL from environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = os.Getenv("PG_DSN")
	}
	if dbURL == "" {
		dbURL = "postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable"
	}

	if *dryRun {
		log.Println("DRY RUN MODE - No changes will be made")
		log.Printf("Would seed database: %s", maskPassword(dbURL))
		log.Printf("Tenant ID: %s", *tenant)
		log.Printf("Clean first: %v", *clean)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Connect to database
	log.Println("Connecting to database...")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Create seeder
	seeder := seed.NewSeeder(pool, *tenant, *verbose)

	// Clean if requested
	if *clean {
		log.Println("Cleaning existing data...")
		if err := seeder.Clean(ctx); err != nil {
			log.Fatalf("Failed to clean database: %v", err)
		}
	}

	// Run seeder
	log.Println("Seeding database...")
	if err := seeder.Run(ctx); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Done!")
}

// maskPassword masks the password in a database URL for logging.
func maskPassword(url string) string {
	// Simple masking - find :password@ pattern and replace password
	result := ""
	inPassword := false
	afterColon := false

	for i, c := range url {
		if c == ':' && i > 10 { // Skip protocol colon
			afterColon = true
			result += string(c)
			continue
		}
		if afterColon && c == '@' {
			inPassword = false
			afterColon = false
			result += string(c)
			continue
		}
		if afterColon && !inPassword {
			inPassword = true
			result += "****"
		}
		if !inPassword {
			result += string(c)
		}
	}

	if result == "" {
		return url
	}
	return result
}

func init() {
	// Configure log format
	log.SetFlags(log.Ltime)
	log.SetPrefix("[seed] ")

	// Print usage if help requested
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: seed [options]\n\n")
		fmt.Fprintf(os.Stderr, "Seeds the development database with realistic Kenyan test data.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment variables:\n")
		fmt.Fprintf(os.Stderr, "  DB_URL or PG_DSN    PostgreSQL connection string\n")
		fmt.Fprintf(os.Stderr, "                      (default: postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  seed                    # Seed with default tenant\n")
		fmt.Fprintf(os.Stderr, "  seed --clean            # Clean and re-seed\n")
		fmt.Fprintf(os.Stderr, "  seed --tenant=my-tenant # Seed for specific tenant\n")
		fmt.Fprintf(os.Stderr, "  seed --dry-run          # Show what would happen\n")
	}
}
