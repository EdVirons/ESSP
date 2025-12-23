package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Minimal migrator for the skeleton.
// - Executes *.sql files in ./migrations in lexical order
// - If a file contains Goose markers, only the `-- +goose Up` section is executed.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: migrate up")
		os.Exit(2)
	}
	mode := os.Args[1]
	if mode != "up" {
		fmt.Println("only 'up' is supported in this skeleton migrator")
		os.Exit(2)
	}

	dir := filepath.Join(".", "migrations")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	files, _ := filepath.Glob(filepath.Join(dir, "*.sql"))
	sort.Strings(files)

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			panic(err)
		}

		sql := extractGooseUp(string(b))
		sql = strings.TrimSpace(sql)
		if sql == "" {
			fmt.Println("skipped", filepath.Base(f), "(empty up section)")
			continue
		}
		_, err = pool.Exec(context.Background(), sql)
		if err != nil {
			panic(fmt.Errorf("migration %s failed: %w", filepath.Base(f), err))
		}
		fmt.Println("applied", filepath.Base(f))
	}
}

func extractGooseUp(in string) string {
	upMarker := "-- +goose Up"
	downMarker := "-- +goose Down"

	up := strings.Index(in, upMarker)
	if up == -1 {
		return in
	}
	start := up + len(upMarker)
	down := strings.Index(in, downMarker)
	if down == -1 {
		return in[start:]
	}
	return in[start:down]
}
