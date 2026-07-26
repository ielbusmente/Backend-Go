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
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("No .env file found, relying on system environment variables")
	// }

	// // Fetch the connection string from the environment variable
	// connString := os.Getenv("SUPABASE_URL")
	// if connString == "" {
	// 	log.Fatal("SUPABASE_URL environment variable is not set")
	// }

	// // Connect to the database
	// ctx := context.Background()
	// conn, err := pgx.Connect(ctx, connString)
	// if err != nil {
	// 	log.Fatalf("Unable to connect to database: %v\n", err)
	// }
	// defer conn.Close(ctx)

	// // Test the connection
	// var greeting string
	// err = conn.QueryRow(ctx, "SELECT 'Hello from Supabase & Go via Env!'").Scan(&greeting)
	// if err != nil {
	// 	log.Fatalf("Query failed: %v\n", err)
	// }

	// fmt.Println(greeting)
	app := fiber.New()
	app.Use(cors.New())

	db, err := initDB()

	if err != nil {
		log.Fatalln("DB Connection Failed")
	}

	defer db.Close()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/records", func(c fiber.Ctx) error {
		ctx := c.Context()

		var count int
		// Query using Fiber's request context
		err := db.QueryRow(ctx, `SELECT COUNT(*) FROM golangtest`).Scan(&count)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch records count" + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"total_records": count,
		})
	})
	// app.Get("/records", func(c fiber.Ctx) error {
	// 	return getAllRecords(c, db)
	// })

	log.Fatal(app.Listen(checkPort()))
}
