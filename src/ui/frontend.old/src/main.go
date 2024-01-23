package main

import (
	"fmt"
	"log"
	"strconv"
	"ui/frontend/src/storage"

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

	//srvAddrs := []string{"localhost:5000", "localhost:5001", "localhost:5002", "localhost:5003"}
	srvAddrs := []string{"node1", "node2", "node3", "node4"}
	nodes := CreateNodes(srvAddrs)

	client := storage.NewStorageClient(nodes.GetAddresses())
	nodes.AddIDs(client.GetNodesInfo())

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("ui", fiber.Map{
			"PageTitle":      "UI",
			"UpdateInterval": 5,
			"Nodes":          nodes,
			"NumNodes":       len(nodes),
		}, "layouts/StaticLayout")
	})

	app.Get("/register/:id", func(c *fiber.Ctx) error {
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
		node := nodes.GetNode(uint32(id))
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
		val, _ := client.ReadSingle(uint32(id))
		return c.JSON(val)
	})

	app.Get("/read", func(c *fiber.Ctx) error {
		value, err := client.ReadValue()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(value)
		return c.JSON(value)
	})

	app.Get("/write/:value", func(c *fiber.Ctx) error {
		value := c.Params("value")
		err := client.WriteValue(value)
		if err != nil {
			log.Println(err)
		}
		return c.JSON(value)
	})

	log.Fatal(app.Listen(":80"))
}
