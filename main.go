package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Load the .env file (optional in production, but great for local development)
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

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
