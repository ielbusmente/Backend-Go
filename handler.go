package main

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
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

func getRecordByID(c fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing ID parameter",
		})
	}
	var record Record
	var query = `
		SELECT id, name, details, created_at 
		FROM golangtest 
		WHERE id=$1
	`
	err := db.QueryRow(c, query, id).Scan(&record.ID, &record.Name, &record.Details, &record.CreatedAt)
	if err != nil {
		return c.Status(404).SendString("Record not found")
	}

	return c.JSON(record)
}

func createRecord(c fiber.Ctx, db *pgxpool.Pool) error {
	var record Record

	if err := c.Bind().Body(&record); err != nil {
		return c.Status(400).SendString("Invalid request body: " + err.Error())
	}

	// Insert the new record into the database
	var query = `
		INSERT INTO golangtest (name, details) 
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := db.QueryRow(c, query, record.Name, record.Details).Scan(
		&record.ID,
		&record.CreatedAt,
	)
	if err != nil {
		return c.Status(500).SendString("Failed to insert record: " + err.Error())
	}
	return c.Status(201).JSON(fiber.Map{
		"record": record,
	})
}

func updateRecord(c fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing ID parameter",
		})
	}
	var record Record

	if err := c.Bind().Body(&record); err != nil {
		return c.Status(400).SendString("Invalid request body: " + err.Error())
	}

	// Update the record in the database
	var query = `
		UPDATE golangtest 
		SET name=$1, details=$2 
		WHERE id=$3
	`
	result, err := db.Exec(c, query, record.Name, record.Details, id)
	if err != nil {
		return c.Status(500).SendString("Failed to update record: " + err.Error())
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(404).SendString("Record not found")
	}

	return c.JSON(record)
}

func deleteRecord(c fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing ID parameter",
		})
	}

	var name string
	// Delete the record from the database
	var query = `
		DELETE FROM golangtest 
		WHERE id=$1
		RETURNING name
	`

	err := db.QueryRow(c, query, id).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(404).JSON(fiber.Map{
				"error": fmt.Sprintf("Record with ID '%s' not found", id),
			})
		}

		return c.Status(500).SendString("Failed to delete record: " + err.Error())
	}

	return c.Status(200).JSON(fiber.Map{
		"message": fmt.Sprintf("Record `%s` deleted successfully", name),
		"id":      id,
	})
}
