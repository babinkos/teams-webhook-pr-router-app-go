package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/romana/rlog"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	os.Setenv("RLOG_LOG_STREAM", "stdout")
	rlog.UpdateEnv()
	var logLevel string = os.Getenv("RLOG_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	rlog.Info("RLOG_LOG_LEVEL: ", logLevel)

	app.Use(requestid.New(requestid.Config{
		Next:       nil,
		Header:     fiber.HeaderXRequestID,
		Generator:  utils.UUIDv4,
		ContextKey: "requestid",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] [${ip}]:${port} ${locals:requestid} ${status} - ${latency} ${bytesReceived} ${method} ${path}\n",
	}))

	app.Use(compress.New(compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/healthz"
		},
		Level: compress.LevelBestSpeed, // 1
	}))

	type SomeStruct struct {
		RequestID string
	}

	// GET /healthz
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	// POST /webhookb2/uid1@uid2/IncomingWebhook/uid3/uid4
	app.All("/webhookb2/:id1/IncomingWebhook/:id2/:id3", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "application/json"
		c.AcceptsEncodings("compress", "br")
		data := SomeStruct{
			RequestID: c.GetRespHeader("X-Request-Id"),
		}

		if logLevel != "DEBUG" {
			// https://docs.gofiber.io/api/ctx#path :
			// override Path with sha256 encoded webhook credentials
			id1 := fmt.Sprintf("%x", sha256.Sum256([]byte(c.Params("id1"))))
			id2 := fmt.Sprintf("%x", sha256.Sum256([]byte(c.Params("id2"))))
			id3 := fmt.Sprintf("%x", sha256.Sum256([]byte(c.Params("id3"))))
			newPath := fmt.Sprintf("/webhookb2/%s/IncomingWebhook/%s/%s", id1[0:7], id2[0:7], id3[0:7])
			c.Path(newPath)
		}

		rlog.Debugf("hook ids: %s, %s, %s ; body: %s", c.Params("id1"), c.Params("id2"), c.Params("id3"), c.Body())

		if reflect.DeepEqual(c.Body(), []byte("{\"test\": true}")) {
			rlog.Debug("Request was Test ping ")
			c.Set("Content-Type", "text/plain")
			return c.SendString("ok")
		} else {
			c.Set("Content-Type", "application/json")
			return c.JSON(data)
		}
		// return c.SendStatus(200)
	})

	go func() {
		log.Fatal(app.Listen(":8080"))
	}()

	appHealth := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	// GET /healthz
	appHealth.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})
	log.Fatal(appHealth.Listen(":9000"))

}
