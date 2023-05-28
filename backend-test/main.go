package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// app.Use(requestid.New())
	// app.Use(logger.New(logger.Config{
	// 	// For more options, see the Config section
	// 	Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	// }))

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(compress.New(compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/healthz"
		},
		Level: compress.LevelBestSpeed, // 1
	}))

	// GET /healthz
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})
	// type SomeStruct struct {
	// 	HookIDs string
	// 	// OriginalBody string
	// }

	// POST /webhookb2/uid1@uid2/IncomingWebhook/uid3/uid4
	app.All("/webhookb2/:id1/IncomingWebhook/:id2/:id3", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "application/json"
		c.AcceptsEncodings("compress", "br")
		// data := SomeStruct{
		// 	HookIDs: fmt.Sprintf("%s, %s, %s", c.Params("id1"), c.Params("id2"), c.Params("id3")),
		// 	OriginalBody: fmt.Sprintf("%s", c.Body()),
		// }
		fmt.Println(fmt.Sprintf("DEBUG: hook ids: %s, %s, %s ; body: %s", c.Params("id1"), c.Params("id2"), c.Params("id3"), c.Body()))
		// return c.JSON(data)
		return c.SendStatus(200)
	})

	log.Fatal(app.Listen(":3000"))
}
