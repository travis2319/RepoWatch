// internal/api/routes.go
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/travis2319/RepoWatch/internal/services"
)

func SetupRoutes(app *fiber.App, checkerService services.CheckerServiceInterface) {
	// Add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API routes
	api := app.Group("/api/v1")

	api.Post("/check", func(c *fiber.Ctx) error {
		// Define request body structure
		type CheckRequest struct {
			Owner string `json:"owner" validate:"required"`
			User  string `json:"user" validate:"required"`
		}

		var req CheckRequest

		// Parse JSON body
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		}

		// Validate required fields
		if req.Owner == "" || req.User == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "owner and user fields are required",
			})
		}

		results, err := checkerService.CheckCollaborators(req.Owner, req.User)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data":    results,
			"message": "Successfully checked collaborators",
		})
	})

	api.Post("/check-single", func(c *fiber.Ctx) error {
		type CheckSingleRequest struct {
			Owner string `json:"owner" validate:"required"`
			Repo  string `json:"repo" validate:"required"`
			User  string `json:"user" validate:"required"`
		}

		var req CheckSingleRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if req.Owner == "" || req.Repo == "" || req.User == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "owner, repo, and user are required fields",
			})
		}

		result, err := checkerService.CheckSingleRepo(req.Owner, req.Repo, req.User)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data":    result,
			"message": "Successfully checked single repository",
		})
	})

	// Legacy endpoint for backwards compatibility
	app.Get("/check-single", func(c *fiber.Ctx) error {
		owner := "travis2319"
		repo := "PROJECT-ADAM"
		// user := "VOID-001"
		user := "ChetanNaikk"

		result, err := checkerService.CheckSingleRepo(owner, repo, user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(result)
	})

	api.Get("/repos", func(c *fiber.Ctx) error {
    owner := c.Query("owner")
    if owner == "" {
        return c.Status(400).JSON(fiber.Map{"error": "owner query param required"})
    }
    repos, err := checkerService.GetOwnerRepos(owner)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"data": repos, "message": "repos fetched"})
})
}
