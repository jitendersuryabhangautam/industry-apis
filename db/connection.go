// Package db handles database connection initialization and management.
// It provides global database pool instance for application-wide access.
package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global database connection pool used throughout the application.
// It is initialized in the Init() function and should be closed using Close() on shutdown.
var DB *pgxpool.Pool

// Init initializes the database connection pool by reading the connection URL from environment variables.
// It creates a connection pool configuration, establishes the connection, and verifies connectivity.
// Returns an error if:
// - DB_URL environment variable is not set
// - Connection string parsing fails
// - Connection pool creation fails
// - Database ping fails (connectivity check)
func Init() error {
	// Get the database connection URL from environment variables
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return fmt.Errorf("DB_URL not set in environment variables")
	}
	// Parse the connection string to create a pool configuration
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %v", err)
	}
	// Create a new connection pool with the parsed configuration
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %v", err)
	}
	// Verify database connectivity by sending a ping
	if err = DB.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	return nil
}

// Close gracefully closes the database connection pool.
// This should be called when the application is shutting down to release database resources.
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("✅ Database connection closed")
	}
}
