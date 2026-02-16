package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	// Import the go-libsql driver for Turso/SQLite
	_ "github.com/tursodatabase/go-libsql"
)

// NewDBConnection initializes and returns a new database connection.
func NewDBConnection(databaseURL string) (*sql.DB, error) {
	driverName := "libsql"
	var dsn string

	if strings.HasPrefix(databaseURL, "libsql://") {
		dsn = databaseURL
		log.Printf("Connecting to Turso database (via libsql driver): %s", databaseURL)
	} else {
		// Assume it's a local SQLite file path
		dbPath := databaseURL

		// Conditionally create the local SQLite database file if it doesn't exist
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			log.Printf("Local SQLite database file does not exist, creating %s", dbPath)
			file, err := os.Create(dbPath)
			if err != nil {
				return nil, fmt.Errorf("failed to create local database file: %w", err)
			}
			file.Close()
			log.Printf("Local SQLite database file created at %s", dbPath)
		}
		// The `_foreign_keys=on` parameter ensures foreign key constraints are enforced for SQLite.
		dsn = fmt.Sprintf("file:%s?_foreign_keys=on", dbPath)
		log.Printf("Connecting to local SQLite database (via libsql driver): %s", dbPath)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection with driver %s: %w", driverName, err)
	}

	// Ping the database to verify the connection is alive.
	if err = db.Ping(); err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("failed to connect to the database with driver %s: %w", driverName, err)
	}

	log.Printf("Successfully connected to database using driver %s", driverName)
	return db, nil
}
