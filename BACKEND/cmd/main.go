// cmd/main.go
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/travis2319/GITHUB-ACCESS/internal/api"
	"github.com/travis2319/GITHUB-ACCESS/internal/github"
	"github.com/travis2319/GITHUB-ACCESS/internal/repository"
	"github.com/travis2319/GITHUB-ACCESS/internal/services"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system env variables")
	}
	// Check for GitHub token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("⚠️  Warning: GITHUB_TOKEN not set! API calls will fail.")
	} else {
		log.Println("✅ GITHUB_TOKEN is set.")
	}

	// Initialize Fiber app
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

	// Setup routes
	api.SetupRoutes(app, checkerService)

	log.Printf("🚀 Server starting on port 4000...")
	log.Fatal(app.Listen(":4000"))
}
