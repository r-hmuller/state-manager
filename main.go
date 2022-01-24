package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()

	app.Get("/:service", func(c *fiber.Ctx) error {
		service := c.Params("service")
		return c.SendString(service)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		rawPayload := c.Body()
		var payload Registry
		json.Unmarshal(rawPayload, &payload)

		return c.SendString("Criado com sucesso")
	})
	log.Fatal(app.Listen(":3000"))
}
