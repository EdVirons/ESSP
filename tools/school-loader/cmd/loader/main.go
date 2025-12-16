// Package main provides a CLI tool to load school data from the scrapper SQLite database
// and import it into the ssot-school service.
package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// County represents a county from the scrapper database
type County struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// SubCounty represents a sub-county from the scrapper database
type SubCounty struct {
	ID       string `json:"id"`
	CountyID string `json:"countyId"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}

// School represents a school from the scrapper database
type School struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Code          string  `json:"code"`
	CountyID      string  `json:"countyId"`
	SubCountyID   string  `json:"subCountyId"`
	Level         string  `json:"level"`
	Type          string  `json:"type"`
	Active        bool    `json:"active"`
	KnecCode      string  `json:"knecCode,omitempty"`
	Uic           string  `json:"uic,omitempty"`
	Sex           string  `json:"sex,omitempty"`
	Cluster       string  `json:"cluster,omitempty"`
	Accommodation string  `json:"accommodation,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
}

// ImportPayload is the payload sent to the ssot-school service
type ImportPayload struct {
	Counties    []County    `json:"counties"`
	SubCounties []SubCounty `json:"subCounties"`
	Schools     []School    `json:"schools"`
}

// ImportResponse is the response from the ssot-school service
type ImportResponse struct {
	OK       bool           `json:"ok"`
	Imported map[string]int `json:"imported"`
}

// Config holds the CLI configuration
type Config struct {
	SQLitePath string
	SSOTURL    string
	TenantID   string
	BatchSize  int
	DryRun     bool
	Verbose    bool
}

func main() {
	cfg := Config{}

	flag.StringVar(&cfg.SQLitePath, "sqlite-path", "/home/pato/scrapper/data/schools.db", "Path to the SQLite database")
	flag.StringVar(&cfg.SSOTURL, "ssot-url", "http://localhost:8081", "URL of the ssot-school service")
	flag.StringVar(&cfg.TenantID, "tenant-id", "tenant_001", "Tenant ID for the import")
	flag.IntVar(&cfg.BatchSize, "batch-size", 1000, "Number of records per batch")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "If true, do not send data to ssot-school")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	if cfg.SQLitePath == "" {
		log.Fatal("--sqlite-path is required")
	}
	if cfg.TenantID == "" {
		log.Fatal("--tenant-id is required")
	}

	if err := run(cfg); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(cfg Config) error {
	// Open SQLite database
	db, err := sql.Open("sqlite", cfg.SQLitePath+"?mode=ro")
	if err != nil {
		return fmt.Errorf("failed to open SQLite database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to SQLite database: %w", err)
	}
	log.Printf("Connected to SQLite database: %s", cfg.SQLitePath)

	// Load all data
	counties, err := loadCounties(db, cfg.Verbose)
	if err != nil {
		return fmt.Errorf("failed to load counties: %w", err)
	}
	log.Printf("Loaded %d counties", len(counties))

	subCounties, err := loadSubCounties(db, cfg.Verbose)
	if err != nil {
		return fmt.Errorf("failed to load sub-counties: %w", err)
	}
	log.Printf("Loaded %d sub-counties", len(subCounties))

	schools, err := loadSchools(db, cfg.Verbose)
	if err != nil {
		return fmt.Errorf("failed to load schools: %w", err)
	}
	log.Printf("Loaded %d schools", len(schools))

	if cfg.DryRun {
		log.Println("Dry run mode - not sending data to ssot-school")
		printSample(counties, subCounties, schools)
		return nil
	}

	// Import in batches
	return importData(cfg, counties, subCounties, schools)
}

func loadCounties(db *sql.DB, verbose bool) ([]County, error) {
	rows, err := db.Query(`
		SELECT id, name, COALESCE(code, '') as code
		FROM counties
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counties []County
	for rows.Next() {
		var c County
		var id int64
		if err := rows.Scan(&id, &c.Name, &c.Code); err != nil {
			return nil, err
		}
		// Generate a stable ID from the source
		c.ID = fmt.Sprintf("county_%d", id)
		if c.Code == "" {
			c.Code = fmt.Sprintf("%03d", id)
		}
		counties = append(counties, c)
		if verbose {
			log.Printf("  County: %s (%s)", c.Name, c.Code)
		}
	}
	return counties, rows.Err()
}

func loadSubCounties(db *sql.DB, verbose bool) ([]SubCounty, error) {
	rows, err := db.Query(`
		SELECT sc.id, sc.county_id, sc.name, COALESCE(sc.code, '') as code
		FROM sub_counties sc
		ORDER BY sc.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subCounties []SubCounty
	for rows.Next() {
		var sc SubCounty
		var id, countyID int64
		if err := rows.Scan(&id, &countyID, &sc.Name, &sc.Code); err != nil {
			return nil, err
		}
		sc.ID = fmt.Sprintf("subcounty_%d", id)
		sc.CountyID = fmt.Sprintf("county_%d", countyID)
		if sc.Code == "" {
			sc.Code = fmt.Sprintf("%05d", id)
		}
		subCounties = append(subCounties, sc)
		if verbose {
			log.Printf("  SubCounty: %s (%s) -> County %s", sc.Name, sc.Code, sc.CountyID)
		}
	}
	return subCounties, rows.Err()
}

func loadSchools(db *sql.DB, verbose bool) ([]School, error) {
	// Query schools with JOIN to get county_id through sub_counties
	query := `
		SELECT
			s.id,
			s.name,
			sc.county_id,
			s.sub_county_id,
			COALESCE(s.knec_code, '') as knec_code,
			COALESCE(s.uic, '') as uic,
			COALESCE(s.level, 'Other') as level,
			COALESCE(s.school_type, 'public') as school_type,
			COALESCE(s.sex, '') as sex,
			COALESCE(s.cluster, '') as cluster,
			COALESCE(s.accommodation, '') as accommodation,
			COALESCE(s.latitude, '0') as latitude,
			COALESCE(s.longitude, '0') as longitude
		FROM schools s
		JOIN sub_counties sc ON s.sub_county_id = sc.id
		ORDER BY s.name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []School
	for rows.Next() {
		var s School
		var id, countyID, subCountyID int64
		var knecCode, uic, level, typ, sex, cluster, accommodation string
		var latStr, lngStr string

		if err := rows.Scan(&id, &s.Name, &countyID, &subCountyID,
			&knecCode, &uic, &level, &typ, &sex, &cluster, &accommodation, &latStr, &lngStr); err != nil {
			return nil, err
		}

		s.ID = fmt.Sprintf("school_%d", id)
		s.CountyID = fmt.Sprintf("county_%d", countyID)
		s.SubCountyID = fmt.Sprintf("subcounty_%d", subCountyID)
		s.Code = knecCode // Use knec_code as the primary code
		s.KnecCode = knecCode
		s.Uic = uic
		s.Level = normalizeLevel(level)
		s.Type = normalizeType(typ)
		s.Sex = sex
		s.Cluster = cluster
		s.Accommodation = accommodation
		s.Latitude = parseFloat(latStr)
		s.Longitude = parseFloat(lngStr)
		s.Active = true

		schools = append(schools, s)
		if verbose && len(schools) <= 10 {
			log.Printf("  School: %s (%s) Level=%s Type=%s", s.Name, s.Code, s.Level, s.Type)
		}
	}
	return schools, rows.Err()
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func getTableColumns(db *sql.DB, table string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}
	return columns, rows.Err()
}

func normalizeLevel(level string) string {
	level = strings.TrimSpace(strings.ToUpper(level))
	switch level {
	case "JSS":
		return "JSS" // Junior Secondary School
	case "SSS":
		return "SSS" // Senior Secondary School
	case "PRIMARY", "PRI":
		return "Primary"
	case "TVET", "TERTIARY":
		return "TVET"
	case "UNIVERSITY":
		return "University"
	default:
		return level // Keep original if not recognized
	}
}

func normalizeType(typ string) string {
	typ = strings.TrimSpace(strings.ToLower(typ))
	switch typ {
	case "public", "government", "gov":
		return "public"
	case "private":
		return "private"
	default:
		return "public"
	}
}

func importData(cfg Config, counties []County, subCounties []SubCounty, schools []School) error {
	client := &http.Client{Timeout: 60 * time.Second}

	totalCounties := len(counties)
	totalSubCounties := len(subCounties)
	totalSchools := len(schools)

	log.Printf("Starting import: %d counties, %d sub-counties, %d schools", totalCounties, totalSubCounties, totalSchools)

	// First, import counties and sub-counties in one batch (they are small)
	payload := ImportPayload{
		Counties:    counties,
		SubCounties: subCounties,
	}

	log.Println("Importing counties and sub-counties...")
	resp, err := sendImport(client, cfg.SSOTURL, cfg.TenantID, payload)
	if err != nil {
		return fmt.Errorf("failed to import counties/sub-counties: %w", err)
	}
	log.Printf("  Imported: counties=%d, subCounties=%d",
		resp.Imported["counties"], resp.Imported["subCounties"])

	// Import schools in batches
	for i := 0; i < len(schools); i += cfg.BatchSize {
		end := i + cfg.BatchSize
		if end > len(schools) {
			end = len(schools)
		}
		batch := schools[i:end]

		log.Printf("Importing schools batch %d-%d of %d...", i+1, end, totalSchools)

		payload := ImportPayload{
			Schools: batch,
		}

		resp, err := sendImport(client, cfg.SSOTURL, cfg.TenantID, payload)
		if err != nil {
			return fmt.Errorf("failed to import schools batch %d-%d: %w", i+1, end, err)
		}
		log.Printf("  Batch imported: schools=%d", resp.Imported["schools"])
	}

	log.Println("Import complete!")
	return nil
}

func sendImport(client *http.Client, baseURL, tenantID string, payload ImportPayload) (*ImportResponse, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/import", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var result ImportResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func printSample(counties []County, subCounties []SubCounty, schools []School) {
	log.Println("\n=== SAMPLE DATA ===")

	log.Println("\n--- Counties (first 5) ---")
	for i, c := range counties {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s: %s (code=%s)\n", c.ID, c.Name, c.Code)
	}

	log.Println("\n--- Sub-Counties (first 5) ---")
	for i, sc := range subCounties {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s: %s (code=%s, county=%s)\n", sc.ID, sc.Name, sc.Code, sc.CountyID)
	}

	log.Println("\n--- Schools (first 10) ---")
	for i, s := range schools {
		if i >= 10 {
			break
		}
		fmt.Printf("  %s: %s\n", s.ID, s.Name)
		fmt.Printf("    Code=%s, Level=%s, Type=%s\n", s.Code, s.Level, s.Type)
		fmt.Printf("    County=%s, SubCounty=%s\n", s.CountyID, s.SubCountyID)
		if s.KnecCode != "" {
			fmt.Printf("    KNEC=%s, UIC=%s, Sex=%s\n", s.KnecCode, s.Uic, s.Sex)
		}
	}

	// Generate sample JSON payload
	log.Println("\n--- Sample Import Payload (JSON) ---")
	payload := ImportPayload{
		Counties:    counties[:min(3, len(counties))],
		SubCounties: subCounties[:min(3, len(subCounties))],
		Schools:     schools[:min(5, len(schools))],
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	os.Stdout.Write(data)
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
