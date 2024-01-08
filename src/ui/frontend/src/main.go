package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Create a new engine
	engine := html.New("./views", ".html")

	// Or from an embedded system
	// See github.com/gofiber/embed for examples
	// engine := html.NewFileSystem(http.Dir("./views", ".html"))

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")

	nodes := TestNodes()

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("ui", fiber.Map{
			"PageTitle":      "UI",
			"UpdateInterval": 5,
			"Nodes":          nodes,
			"NumNodes":       len(nodes),
		}, "layouts/StaticLayout")
	})

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(nodes)
	})

	log.Fatal(app.Listen(":80"))
}
