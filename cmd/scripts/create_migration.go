package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 2 {
		fmt.Println("	Usage: go run create_migration.go <name>")
		os.Exit(1)
	}

	name := strings.ToLower(os.Args[1])

	migrationsDir := "schema/migrations"
	if err := os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
		fmt.Printf("Failed to create migrations dir: %v\n", err)
		os.Exit(1)
	}

	index, err := getNextMigrationIndex(migrationsDir)
	if err != nil {
		fmt.Printf("Failed to get next migration index: %v\n", err)
		os.Exit(1)
	}

	version := fmt.Sprintf("%06d", index)
	base := fmt.Sprintf("%s_%s", version, name)
	upPath := filepath.Join(migrationsDir, base+".up.sql")
	downPath := filepath.Join(migrationsDir, base+".down.sql")

	upSQL, downSQL := generateTemplates()

	if err := os.WriteFile(upPath, []byte(upSQL), 0644); err != nil {
		fmt.Printf("Failed to write %s: %v\n", upPath, err)
		os.Exit(1)
	}
	if err := os.WriteFile(downPath, []byte(downSQL), 0644); err != nil {
		fmt.Printf("Failed to write %s: %v\n", downPath, err)
		os.Exit(1)
	}

	fmt.Printf("Created:\n - %s\n - %s\n", upPath, downPath)
}

func getNextMigrationIndex(dir string) (int, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 1, nil // first migration
	}

	re := regexp.MustCompile(`^(\d{6})_.*\.up\.sql$`)
	var indices []int

	for _, f := range files {
		match := re.FindStringSubmatch(f.Name())
		if len(match) == 2 {
			n, err := strconv.Atoi(match[1])
			if err == nil {
				indices = append(indices, n)
			}
		}
	}

	if len(indices) == 0 {
		return 1, nil
	}

	sort.Ints(indices)
	return indices[len(indices)-1] + 1, nil
}

func generateTemplates() (string, string) {
	return `CREATE TABLE
    IF NOT EXISTS your_table (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`,
		`DROP TABLE IF EXISTS your_table;`
}
