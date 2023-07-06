package main

import (
	fragments "github.com/katallaxie/fiber-mosaic"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Create a new engine
	engine := html.New("./views", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Add logger
	app.Use(logger.New())

	app.Get("/index", fragments.Template(fragments.Config{}, "index", fiber.Map{}, "layouts/main"))

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
