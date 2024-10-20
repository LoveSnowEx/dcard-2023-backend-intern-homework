package router

import (
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func Setup(app *fiber.App) {
	app.Use(healthcheck.New())
	app.Get("/api/head/:key", handler.GetHead())
	app.Get("/api/page/:key", handler.GetPage())
}
