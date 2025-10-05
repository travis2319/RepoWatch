package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/travis2319/RepoWatch/internal/api"
	"github.com/travis2319/RepoWatch/internal/github"
	"github.com/travis2319/RepoWatch/internal/repository"
	"github.com/travis2319/RepoWatch/internal/services"
)

func main() {
	// Optional: load .env if needed
	// _ = godotenv.Load()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("⚠️  Warning: GITHUB_TOKEN not set! API calls will fail.")
	} else {
		log.Println("✅ GITHUB_TOKEN is set.")
	}

	app := fiber.New(fiber.Config{
		AppName: "GitHub Collaborator Checker API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "An error occurred processing your request",
			})
		},
	})

	// Middleware
	app.Use(logger.New())

	// Initialize database
	db := repository.InitDB("db/collabs.db")
	defer db.Close()

	// Initialize repositories
	repoRepo := repository.NewRepoRepository(db)
	collabRepo := repository.NewCollaboratorRepository(db)

	// Initialize GitHub client
	githubClient := github.NewClient()

	// Initialize services
	checkerService := services.NewCheckerService(repoRepo, collabRepo, githubClient)

	// Setup API routes
	api.SetupRoutes(app, checkerService)

	// 🔥 Serve static frontend build
	app.Static("/", "../FRONTEND/build")

	// SPA fallback (to support client-side routing)
	app.All("*", func(c *fiber.Ctx) error {
		return c.SendFile("../FRONTEND/build/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Printf("🚀 Server running at http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
