package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"time"
)

func main() {
	app := fiber.New()

	app.Get("/:service", func(c *fiber.Ctx) error {
		service := c.Params("service")
		serviceUuid := GetLatestCheckpoint(service)
		return c.SendString(serviceUuid)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		rawPayload := c.Body()
		var payload Registry
		err := json.Unmarshal(rawPayload, &payload)
		if err != nil {
			log.Fatal(err)
		}
		payload.Uuid = uuid.New().String()
		payload.Datetime = time.Now().Unix()
		return c.SendStatus(201)
	})
	log.Fatal(app.Listen(":8889"))
}
