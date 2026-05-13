package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/travis2319/RepoWatch/internal/api"
	"github.com/travis2319/RepoWatch/internal/github"
	"github.com/travis2319/RepoWatch/internal/repository"
	"github.com/travis2319/RepoWatch/internal/services"
)

func main() {

	// Load .env locally (optional in Docker)
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ No local .env file found, using system environment variables")
	}

	// Validate GitHub token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("❌ GITHUB_TOKEN is not set")
	}

	log.Println("✅ GITHUB_TOKEN loaded")

	// App setup
	app := fiber.New(fiber.Config{
		AppName: "GitHub Collaborator Checker API",

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(fiberLogger.New())

	// Database
	db := repository.InitDB("db/collabs.db")
	defer db.Close()

	log.Println("✅ Database initialized")

	// Repositories
	repoRepo := repository.NewRepoRepository(db)
	collabRepo := repository.NewCollaboratorRepository(db)

	// GitHub client
	githubClient := github.NewClient()

	// Services
	checkerService := services.NewCheckerService(
		repoRepo,
		collabRepo,
		githubClient,
	)

	// Routes
	api.SetupRoutes(app, checkerService)

	// Optional frontend static hosting
	// app.Static("/", "./build")
	// app.All("*", func(c *fiber.Ctx) error {
	// 	return c.SendFile("./build/index.html")
	// })

	// Port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4000"
	}

	log.Printf("🚀 Server running on port %s", port)

	log.Fatal(app.Listen(":" + port))
}