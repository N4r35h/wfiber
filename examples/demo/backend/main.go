package main

import (
	"github.com/N4r35h/wfiber/examples/demo/backend/models"
	"github.com/N4r35h/wfiber/wfiber"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := wfiber.New(wfiber.WFiberAppConfig{
		APIPrefix:              "/api",
		FrontendFolder:         "./examples/demo/frontend",
		SwaggerDocFolder:       "./examples/demo/backend/swagger",
		GeneratedAPIClientPath: "/src/api/api.gen.ts",
		GenerateClient:         true,
	})
	app.Static("/swagger", "./examples/demo/backend/swagger")
	app.FApp.Get("/swagger-ui", func(c *fiber.Ctx) error {
		return c.Type("html").SendString(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>wfiber Demo Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: "#swagger-ui",
      });
    </script>
  </body>
</html>
`)
	})

	api := app.Group("/api")
	todos := api.Group("/todos")

	todos.Get("/", nil, []models.Todo{}, func(c *fiber.Ctx) error {
		return c.JSON([]models.Todo{
			{ID: 1, Title: "Ship wfiber demo", Done: false},
		})
	})

	todos.Post("/", models.Todo{}, models.Todo{}, func(c *fiber.Ctx) error {
		var body models.Todo
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.JSON(body)
	})

	_ = app.Listen(":8081")
}
