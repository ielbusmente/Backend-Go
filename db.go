package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func initDB() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Fetch the connection string from the environment variable
	connString := os.Getenv("SUPABASE_URL")
	if connString == "" {
		log.Fatal("SUPABASE_URL environment variable is not set")
	}

	// Connect to the database
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Test the connection
	var greeting string
	err = pool.QueryRow(ctx, "SELECT 'Hello from Supabase & Go via Env!'").Scan(&greeting)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}

	fmt.Println(greeting)
	return pool, nil
}
