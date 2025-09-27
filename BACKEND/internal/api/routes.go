// internal/api/routes.go
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/travis2319/GITHUB-ACCESS/internal/services"
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

	api.Get("/check/:owner/:user", func(c *fiber.Ctx) error {
		owner := c.Params("owner")
		user := c.Params("user")

		if owner == "" || user == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "owner and user parameters are required",
			})
		}

		results, err := checkerService.CheckCollaborators(owner, user)
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
		user := "VOID-001"

		result, err := checkerService.CheckSingleRepo(owner, repo, user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(result)
	})
}
