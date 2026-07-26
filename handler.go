package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Record struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Details   *string `json:"details"`
	CreatedAt string  `json:"created_at"`
}

func getAllRecords(c fiber.Ctx, db *pgxpool.Pool) error {
	// Query using Fiber's request context
	rows, err := db.Query(c, "SELECT id, name, details, created_at FROM golangtest")
	// Fetch all records from the database
	if err != nil {
		return c.Status(500).SendString("Failed to execute query: " + err.Error())
	}
	defer rows.Close()

	var records []Record

	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.ID, &record.Name, &record.Details, &record.CreatedAt); err != nil {
			return c.Status(500).SendString("Failed to scan row: " + err.Error())
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).SendString("Error iterating over rows: " + err.Error())
	}

	return c.JSON(records)
}
