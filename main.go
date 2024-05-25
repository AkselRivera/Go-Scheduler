package main

import (
	"go-scheduler/routers"
	"go-scheduler/utils"
	"go-scheduler/utils/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {
	go scheduler.TaskRunner()

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(cors.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(utils.GetHelloWorld())
	})

	v1.Get("/metrics", monitor.New(monitor.Config{Title: "Go - Chenduler ®"}))

	v1.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	v1Tasks := v1.Group("/tasks")
	routers.TaskRouter(v1Tasks)

	err := app.Listen(":10000")
	if err != nil {
		return
	}
}
