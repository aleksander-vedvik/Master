package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")

	nodes := NewNodes()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("ui", fiber.Map{
			"PageTitle":      "UI",
			"UpdateInterval": 5,
			"Nodes":          nodes,
		}, "layouts/StaticLayout")
	})

	app.Get("/register/:addr", func(c *fiber.Ctx) error {
		addr := c.Params("addr")
		nodes = append(nodes, NewNode(addr, len(nodes)+1))
		return nil
	})

	app.Post("/node", func(c *fiber.Ctx) error {
		payload := Payload{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		nodes.UpdateNode(payload)
		return nil
	})

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(nodes)
	})

	app.Get("/start", func(c *fiber.Ctx) error {
		return c.Render("ui", fiber.Map{
			"PageTitle":      "UI",
			"UpdateInterval": 5,
			"Nodes":          nodes,
		}, "layouts/StaticLayout")
	})

	log.Fatal(app.Listen(":80"))
}
