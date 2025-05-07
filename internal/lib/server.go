package lib

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func NewServer() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(recover.New())
	return app
}
