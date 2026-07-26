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

	app.Get("/records", func(c fiber.Ctx) error {
		return getAllRecords(c, db)
	})

	log.Fatal(app.Listen(checkPort()))
}
