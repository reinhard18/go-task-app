package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"task-management/utils"

	_ "github.com/lib/pq"
)

func NewDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables
	schemaPath := "db/sql/schema.sql"
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = db.Exec(string(schemaBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Create a test user if it doesn't exist
	hashedPassword, _ := utils.HashPassword("testpass")
	_, err = db.Exec(
		"INSERT INTO users (username, password) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING",
		"testuser",
		hashedPassword,
	)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
