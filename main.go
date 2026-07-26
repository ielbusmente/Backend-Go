package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func checkPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	return port
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	db, err := initDB()

	if err != nil {
		log.Fatalln("DB Connection Failed")
	}

	defer db.Close()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, Go!")
	})

	app.Post("/records", func(c fiber.Ctx) error {
		// body structure
		// {name: string, details: string}
		return createRecord(c, db)
	})

	app.Get("/records", func(c fiber.Ctx) error {
		return getAllRecords(c, db)
	})

	app.Get("/records/:id", func(c fiber.Ctx) error {
		return getRecordByID(c, db)
	})

	app.Put("/records/:id", func(c fiber.Ctx) error {
		return updateRecord(c, db)
	})

	app.Delete("/records/:id", func(c fiber.Ctx) error {
		return deleteRecord(c, db)
	})

	log.Fatal(app.Listen(checkPort()))
}
