package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Create a new engine
	engine := html.New(".", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Add logger
	app.Use(logger.New())

	// Add compression
	app.Use(compress.New())

	// Add CORS
	app.Use(cors.New())

	app.Get("/fragment1", func(c *fiber.Ctx) error {
		c.Links("https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css", "stylesheet")

		return c.Render("fragment1", fiber.Map{
			"Title": "Example 1",
		})
	})

	app.Get("/fragment2", func(c *fiber.Ctx) error {
		c.Links("https://unpkg.com/react-dom@17/umd/react-dom.development.js", "script", "")

		c.Response().SetStatusCode(403)

		return c.Render("fragment2", fiber.Map{
			"Title": "Example 2",
		})
	})

	app.Get("/fragment3", func(c *fiber.Ctx) error {
		return c.Send(nil)
	})

	app.Get("/fragment4", func(c *fiber.Ctx) error {
		return c.Render("fragment4", fiber.Map{
			"Title": "Example 4",
		})
	})

	app.Get("/fragment5", func(c *fiber.Ctx) error {
		return c.Render("fragment5", fiber.Map{
			"Title": "Example 5",
		})
	})

	app.Get("/fragment6", func(c *fiber.Ctx) error {
		return c.Render("fragment6", fiber.Map{
			"Title": "Example 6",
		})
	})

	app.Get("/fallback", func(c *fiber.Ctx) error {
		return c.Render("fallback", fiber.Map{
			"Title": "Fallback",
		})
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
