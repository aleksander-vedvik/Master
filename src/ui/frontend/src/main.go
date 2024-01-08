package main

import (
	"log"
	"strconv"

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

	app.Get("/node/:id", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		id, _ := strconv.Atoi(c.Params("id"))
		node := GetNode(id, nodes)
		return c.Render("node", fiber.Map{
			"PageTitle": "Node" + c.Params("id"),
			"Id":        id,
			"Status":    node.Status,
			"Messages":  node.Messages,
			"Round":     node.Round,
		}, "layouts/StaticLayoutNode")
	})

	app.Get("/status/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		node := GetNode(id, nodes)
		return c.JSON(node)
	})

	log.Fatal(app.Listen(":80"))
}
