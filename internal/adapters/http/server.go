package http

import "github.com/gofiber/fiber/v2"
import "github.com/gofiber/fiber/v2/middleware/recover"

type FiberServer struct {
	*fiber.App
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "text-to-api",
			AppName:      "text-to-api",
			BodyLimit:    2 * 1024 * 1024, // 2MB
		}),
	}

	// Add recover middleware to the server
	server.App.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
		},
	))

	return server
}
