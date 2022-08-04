package app

import (
	"github.com/todo-list/internal/adapter/inbound/taskhdl"
	"github.com/todo-list/internal/adapter/outbound/taskrps"
	"github.com/todo-list/internal/core/services/tasksvc"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handlers struct {
	Postgres *gorm.DB
	R        *fiber.App
	Logger   *zap.Logger
}

func (h *Handlers) SetupRouter() {

	// initialize Repository
	taskRepo := taskrps.NewTaskPostgres(h.Postgres)

	// initialize Service
	taskService := tasksvc.NewTaskService(h.Logger, taskRepo)

	// initialize Handler
	taskhdl.NewTaskHandler(h.R, taskService)
}
