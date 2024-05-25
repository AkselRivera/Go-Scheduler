package controllers

import (
	"encoding/json"
	"fmt"
	"go-scheduler/models"
	"go-scheduler/utils/jobs"
	"go-scheduler/utils/scheduler"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type TaskRequest struct {
	ExecutionTime string `json:"execution_date"`
	JobID         string `json:"job_id"`
}

func GetTasksController(c *fiber.Ctx) error {
	tasks := scheduler.GetAllTasks()
	return c.JSON(tasks)
}

func CreateTaskController(c *fiber.Ctx) error {

	var request TaskRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	executionTime, err := time.Parse(time.RFC3339, request.ExecutionTime)

	now := time.Now().UTC()

	if err != nil || executionTime.Before(now) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid date format. Use 2024-05-22T19:04:20Z or execution date is in the past",
			"date":  executionTime,
		})
	}

	task := &models.Task{
		ID:            request.JobID,
		ExecutionTime: executionTime,
		Action: func() {
			jobs.InitialJob(request.JobID)
		},
	}
	taskID, err := scheduler.AddTask(task)
	if err != nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "Task with the same ID already exists",
		})
	}
	scheduler.TaskChannel <- task // Enviar tarea al canal fuera de la función AddTask

	return c.JSON(fiber.Map{
		"id": taskID,
	})
}

func CreateTaskBulkController(c *fiber.Ctx) error {

	var request []TaskRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid input",
			"message": err.Error(),
		})
	}

	for index, task := range request {
		executionTime, err := time.Parse(time.RFC3339, task.ExecutionTime)

		now := time.Now().UTC()

		if err != nil || executionTime.Before(now) {
			return c.Status(400).JSON(fiber.Map{
				"error":     "Invalid date format. Use 2024-05-22T19:04:20Z or execution date is in the past",
				"date":      executionTime,
				"job_index": index,
			})
		}

	}

	for _, job := range request {

		executionTime, _ := time.Parse(time.RFC3339, job.ExecutionTime)

		task := &models.Task{
			ID:            job.JobID,
			ExecutionTime: executionTime,
			Action: func() {
				jobs.InitialJob(job.JobID)
			},
		}
		_, err := scheduler.AddTask(task)
		if err != nil {
			log.Error("Job #%s already exists or: %s", task.ID, err.Error())
		}
		scheduler.TaskChannel <- task // Enviar tarea al canal fuera de la función AddTask
	}

	return c.JSON(fiber.Map{
		"message": "Tasks created successfully",
	})
}

func CancelTaskBulkController(c *fiber.Ctx) error {
	var request []string
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid input must be an array of strings",
			"message": err.Error(),
		})
	}

	for _, id := range request {
		scheduler.CancelTask(id)
	}

	return c.JSON(fiber.Map{
		"status": "Tasks added to the cancel queue",
	})
}

func CancelTaskController(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	scheduler.CancelTask(id)
	return c.JSON(fiber.Map{
		"status": "Task cancelled",
	})
}

func ExecutionTaskController(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	data, err := jobs.GetJob(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Job not found",
		})
	}

	batutaResp, err := jobs.ExecCommand(data)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Batuta Responded with error: " + err.Error(),
		})
	}

	jsonBytes, err := json.Marshal(batutaResp)
	if err != nil {
		fmt.Println("Error al convertir el mapa a JSON:", err)
		return c.Status(409).JSON(batutaResp)
	}

	// Convertir los bytes JSON a un string
	jsonString := string(jsonBytes)

	if strings.Contains(jsonString, "RPORT_HOST_DISCONNECTED") {
		return c.Status(409).JSON(batutaResp)
	}

	return c.JSON(batutaResp)

}
