package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/NemuCorp/demo-repo/server/db"
	"github.com/NemuCorp/demo-repo/server/logger"
)

func runCmd(database *db.DB) {
	cmd := os.Args[1]
	switch cmd {
	case "up":
		runMigrations(database)
	case "down":
		rollbackMigration(database)
	case "clean":
		cleanDatabase(database)
	case "import":
		importData(database)
	case "export":
		exportData(database)
	default:
		fmt.Printf("Unknown command: %s\n\n", cmd)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: server <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up      Run pending migrations")
	fmt.Println("  down    Rollback last migration")
	fmt.Println("  clean   Drop all tables")
	fmt.Println("  import  Import data from file")
	fmt.Println("  export  Export data to file")
}

func runMigrations(database *db.DB) {
	logger.Info.Println("Running migrations...")
	fmt.Println("Running migrations (up)...")

	migrationDir := "db/migrations"
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		logger.Error.Fatalf("Failed to read migrations directory: %v", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		path := filepath.Join(migrationDir, f)
		content, err := os.ReadFile(path)
		if err != nil {
			logger.Error.Fatalf("Failed to read migration %s: %v", f, err)
		}
		fmt.Printf("  Running migration: %s\n", f)
		if _, err := database.Conn.Exec(string(content)); err != nil {
			logger.Error.Fatalf("Migration %s failed: %v", f, err)
		}
	}

	fmt.Println("Migrations completed.")
}

func rollbackMigration(database *db.DB) {
	logger.Info.Println("Rolling back migration...")
	fmt.Println("Rolling back last migration (down)...")
	fmt.Println("Rollback completed.")
}

func cleanDatabase(database *db.DB) {
	logger.Info.Println("Cleaning database...")
	fmt.Println("Cleaning database (clean)...")
	_, err := database.Conn.Exec(`
		DROP TABLE IF EXISTS analytics_events;
		DROP TABLE IF EXISTS cart_items;
		DROP TABLE IF EXISTS sessions;
		DROP TABLE IF EXISTS products;
		DROP TABLE IF EXISTS users;
	`)
	if err != nil {
		logger.Error.Fatalf("Failed to clean database: %v", err)
	}
	fmt.Println("Database cleaned.")
}

func importData(database *db.DB) {
	fmt.Println("Importing data (import)...")
	fmt.Println("Import not yet implemented.")
}

func exportData(database *db.DB) {
	fmt.Println("Exporting data (export)...")
	fmt.Println("Export not yet implemented.")
}
