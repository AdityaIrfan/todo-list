package taskhdl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/todo-list/internal/core/domain"
	"github.com/todo-list/internal/core/ports"
	responseErr "github.com/todo-list/internal/error"
	"github.com/todo-list/internal/response"
)

type taskHandler struct {
	app         *fiber.App
	taskService ports.TaskService
}

func NewTaskHandler(app *fiber.App, taskService ports.TaskService) {
	taskHandler := taskHandler{
		app:         app,
		taskService: taskService,
	}

	api := taskHandler.app.Group("/task")
	api.Post("/add", taskHandler.create)
	api.Get("/get/:id", taskHandler.getOneById)
	api.Put("/update/:id", taskHandler.update)
	api.Delete("/delete/:id", taskHandler.delete)
	api.Get("/get", taskHandler.getAllWithPaginate)
}

func (instance *taskHandler) create(c *fiber.Ctx) error {
	request := new(domain.CreateTaskRequst)
	if err := c.BodyParser(&request); err != nil {
		return responseErr.Response(c, responseErr.ResponseBadRequest(err.Error()))
	}

	// validate here

	if err := instance.taskService.Create(c.Context(), request); err != nil {
		return responseErr.Response(c, err)
	}

	return response.Success(c, fiber.StatusOK, response.SuccessMessage("Success"))
}

func (instance *taskHandler) getOneById(c *fiber.Ctx) error {
	task, err := instance.taskService.GetOneByID(c.Context(), c.Params("id"))
	if err != nil {
		return responseErr.Response(c, err)
	}

	return response.Success(c, fiber.StatusOK, response.SuccessData(task))
}

func (instance *taskHandler) update(c *fiber.Ctx) error {
	request := new(domain.UpdateTaskRequest)
	if err := c.BodyParser(&request); err != nil {
		return responseErr.Response(c, responseErr.ResponseBadRequest(err.Error()))
	}

	// validate here ...

	if err := instance.taskService.Update(c.Context(), c.Params("id"), request); err != nil {
		return responseErr.Response(c, err)
	}

	return response.Success(c, fiber.StatusOK, response.SuccessMessage("Success"))
}

func (instance *taskHandler) delete(c *fiber.Ctx) error {
	if err := instance.taskService.Delete(c.Context(), c.Params("id")); err != nil {
		return responseErr.Response(c, err)
	}

	return response.Success(c, fiber.StatusOK, response.SuccessMessage("Success"))
}

func (instance *taskHandler) getAllWithPaginate(c *fiber.Ctx) error {
	params := new(domain.TaskParams)
	if err := c.QueryParser(params); err != nil {
		return responseErr.Response(c, responseErr.ResponseBadRequest(err.Error()))
	}

	if err := params.Validate(); err != nil {
		return responseErr.Response(c, responseErr.ResponseBadRequest(err.Error()))
	}

	tasks, err := instance.taskService.GetAllWithPaginate(c.Context(), params)
	if err != nil {
		return responseErr.Response(c, err)
	}

	return response.Success(c, fiber.StatusOK, response.SuccessData(tasks))
}
