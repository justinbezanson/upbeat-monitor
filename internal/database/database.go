package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Import the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// NewDBConnection initializes and returns a new database connection.
func NewDBConnection(dbPath string) (*sql.DB, error) {
	// Check if the database file exists, create it if it doesn't
	// This ensures that `local.db` exists before `sql.Open` tries to connect to it.
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Database file does not exist, creating %s", dbPath)
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %w", err)
		}
		file.Close()
		log.Printf("Database file created at %s", dbPath)
	}

	// Open the SQLite database connection.
	// The `_foreign_keys=on` parameter ensures foreign key constraints are enforced.
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Ping the database to verify the connection is alive.
	if err = db.Ping(); err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	log.Printf("Successfully connected to database at %s", dbPath)
	return db, nil
}
