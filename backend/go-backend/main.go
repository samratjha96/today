package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"go-backend/pkg/database"
	"go-backend/pkg/github"
	"go-backend/pkg/hackernews"
)

func main() {
	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// GitHub routes with database persistence
	ghHandler := github.NewHandler()
	ghHandler.RegisterRoutes(app)

	// HackerNews routes with database persistence
	hnHandler := hackernews.NewHandler()
	hnHandler.RegisterRoutes(app)

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
