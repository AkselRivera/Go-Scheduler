package routers

import (
	"go-scheduler/controllers"

	"github.com/gofiber/fiber/v2"
)

func TaskRouter(router fiber.Router) {

	router.Get("/", controllers.GetTasksController)
	router.Post("/", controllers.CreateTaskController)
	router.Post("/bulk", controllers.CreateTaskBulkController)
	router.Delete("/bulk", controllers.CancelTaskBulkController)
	router.Delete("/:id", controllers.CancelTaskController)
}
