package main

import (
	"log"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/internal/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		// Close the database connection when the program exits
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}()

	_ = dbConn

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	router.Setup(app)

	app.Listen(":3000")
}
