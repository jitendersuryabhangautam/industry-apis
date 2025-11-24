package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() error {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return fmt.Errorf("DB_URL not set in environment variables")
	}
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %v", err)
	}
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %v", err)
	}
	if err = DB.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("✅ Database connection closed")
	}
}
