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

// ==========================================
// REUSABLE HELPER FUNCTIONS
// ==========================================

// sendError standardizes error responses in JSON format
func sendError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// sendSuccess standardizes successful JSON responses
func sendSuccess(c fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}

// validateID ensures route parameters aren't empty
func validateID(c fiber.Ctx) (string, error) {
	id := c.Params("id")
	if id == "" {
		return "", errors.New("missing ID parameter")
	}
	return id, nil
}

// ==========================================
// CONTROLLERS / HANDLERS
// ==========================================

func getAllRecords(c fiber.Ctx, db *pgxpool.Pool) error {
	rows, err := db.Query(c, "SELECT id, name, details, created_at FROM golangtest")
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to fetch records: "+err.Error())
	}
	defer rows.Close()

	var records []Record

	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.ID, &record.Name, &record.Details, &record.CreatedAt); err != nil {
			return sendError(c, fiber.StatusInternalServerError, "Failed to parse record data: "+err.Error())
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Error processing record stream: "+err.Error())
	}

	return sendSuccess(c, fiber.StatusOK, records)
}

func getRecordByID(c fiber.Ctx, db *pgxpool.Pool) error {
	id, err := validateID(c)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid ID parameter: "+err.Error())
	}

	var record Record
	query := `
		SELECT id, name, details, created_at 
		FROM golangtest 
		WHERE id = $1
	`
	err = db.QueryRow(c, query, id).Scan(&record.ID, &record.Name, &record.Details, &record.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sendError(c, fiber.StatusNotFound, "Record not found: "+err.Error())
		}
		return sendError(c, fiber.StatusInternalServerError, "Database error: "+err.Error())
	}

	return sendSuccess(c, fiber.StatusOK, record)
}

func createRecord(c fiber.Ctx, db *pgxpool.Pool) error {
	var record Record

	if err := c.Bind().Body(&record); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if record.Name == "" {
		return sendError(c, fiber.StatusBadRequest, "Field 'name' is required")
	}

	query := `
		INSERT INTO golangtest (name, details) 
		VALUES ($1, $2)
		RETURNING id, name, details, created_at
	`
	err := db.QueryRow(c, query, record.Name, record.Details).Scan(
		&record.ID,
		&record.Name,
		&record.Details,
		&record.CreatedAt,
	)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to insert record: "+err.Error())
	}

	return sendSuccess(c, fiber.StatusCreated, fiber.Map{
		"message": "Record created successfully",
		"record":  record,
	})
}

func updateRecord(c fiber.Ctx, db *pgxpool.Pool) error {
	id, err := validateID(c)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid ID parameter: "+err.Error())
	}

	var record Record
	if err := c.Bind().Body(&record); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	// Using COALESCE allows partial updates (if a field isn't provided, keep the existing value)
	query := `
		UPDATE golangtest 
		SET name = COALESCE(NULLIF($1, ''), name), 
		    details = COALESCE($2, details) 
		WHERE id = $3
		RETURNING id, name, details, created_at
	`
	err = db.QueryRow(c, query, record.Name, record.Details, id).Scan(
		&record.ID,
		&record.Name,
		&record.Details,
		&record.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sendError(c, fiber.StatusNotFound, "Record not found: "+err.Error())
		}
		return sendError(c, fiber.StatusInternalServerError, "Failed to update record: "+err.Error())
	}

	return sendSuccess(c, fiber.StatusOK, fiber.Map{
		"message": "Record updated successfully",
		"record":  record,
	})
}

func deleteRecord(c fiber.Ctx, db *pgxpool.Pool) error {
	id, err := validateID(c)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid ID parameter: "+err.Error())
	}

	var name string
	query := `
		DELETE FROM golangtest 
		WHERE id = $1
		RETURNING name
	`

	err = db.QueryRow(c, query, id).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sendError(c, fiber.StatusNotFound, fmt.Sprintf("Record with ID '%s' not found", id))
		}
		return sendError(c, fiber.StatusInternalServerError, "Failed to delete record: "+err.Error())
	}

	return sendSuccess(c, fiber.StatusOK, fiber.Map{
		"message": fmt.Sprintf("Record `%s` deleted successfully", name),
		"id":      id,
	})
}
