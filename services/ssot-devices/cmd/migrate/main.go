package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Simple migrator: runs *.sql files in ./migrations in lexical order.
// In production, swap with goose/atlas/liquibase; this is intentionally minimal.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: migrate [up|down]")
		os.Exit(2)
	}
	dir := filepath.Join(".", "migrations")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://ssp:ssp@localhost:5432/ssp_devices?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil { panic(err) }
	defer pool.Close()

	mode := os.Args[1]
	files, _ := filepath.Glob(filepath.Join(dir, "*.sql"))
	if mode == "up" {
		for _, f := range files {
			b, _ := os.ReadFile(f)
			_, err := pool.Exec(context.Background(), string(b))
			if err != nil { panic(err) }
			fmt.Println("applied", filepath.Base(f))
		}
		return
	}
	if mode == "down" {
		fmt.Println("down not implemented in simple migrator; use drop schema in dev")
		return
	}
	fmt.Println("unknown mode:", mode)
	os.Exit(2)
}
