package tasksvc

import (
	"context"
	"github.com/todo-list/internal/core/domain"
	"github.com/todo-list/internal/core/ports"
	responseErr "github.com/todo-list/internal/error"
	"go.uber.org/zap"
	"strconv"
)

var (
	FailedToCreateNewTask = "Failed to create new task"
	FailedToGetTask       = "Failed to get task"
	FailedToUpdateTask    = "Failed to update task"
	TaskNotFound          = "Task not found"
	FailedToDeleteTask    = "Failed to delete task"
)

type taskService struct {
	log      *zap.Logger
	taskRepo ports.TaskRepository
}

func NewTaskService(log *zap.Logger, taskRepo ports.TaskRepository) ports.TaskService {
	return &taskService{
		log:      log,
		taskRepo: taskRepo,
	}
}

func (instance *taskService) Create(ctx context.Context, request *domain.CreateTaskRequst) error {
	if err := instance.taskRepo.Create(ctx, request.ToBase()); err != nil {
		instance.log.Error("failed to create task : ", zap.Error(err))
		return responseErr.ResponseInternalServerError(FailedToCreateNewTask)
	}

	return nil
}

func (instance *taskService) Update(ctx context.Context, id string, request *domain.UpdateTaskRequest) error {
	// get one by id for checking
	_, err := strconv.Atoi(id)
	if id == "" || err != nil {
		return responseErr.ResponseBadRequest(responseErr.ErrBadRequest.Error())
	}

	task, err := instance.taskRepo.GetOneByID(ctx, id)
	if err != nil {
		instance.log.Error("failed to get task by id ["+id+"] : ", zap.Error(err))
		return responseErr.ResponseInternalServerError(FailedToGetTask)
	}

	if task == nil {
		return responseErr.ResponseNotFound(TaskNotFound)
	}

	if err := instance.taskRepo.Update(ctx, request.ToBase(task)); err != nil {
		instance.log.Error("failed to update task by id ["+id+"]", zap.Error(err))
		return responseErr.ResponseInternalServerError(FailedToUpdateTask)
	}

	return nil
}

func (instance *taskService) GetOneByID(ctx context.Context, id string) (*domain.TaskTransformer, error) {
	_, err := strconv.Atoi(id)
	if id == "" || err != nil {
		return nil, responseErr.ResponseBadRequest(responseErr.ErrBadRequest.Error())
	}

	task, err := instance.taskRepo.GetOneByID(ctx, id)
	if err != nil {
		instance.log.Error("failed to get task by id ["+id+"] : ", zap.Error(err))
		return nil, responseErr.ResponseInternalServerError(FailedToGetTask)
	}

	if task == nil {
		return nil, responseErr.ResponseNotFound(TaskNotFound)
	}

	return task.ToTaskTransformer(), nil
}

func (instance *taskService) Delete(ctx context.Context, id string) error {
	_, err := strconv.Atoi(id)
	if id == "" || err != nil {
		return responseErr.ResponseBadRequest(responseErr.ErrBadRequest.Error())
	}

	task, err := instance.taskRepo.GetOneByID(ctx, id)
	if err != nil {
		instance.log.Error("failed to get task by id ["+id+"] : ", zap.Error(err))
		return responseErr.ResponseInternalServerError(FailedToGetTask)
	}

	if task == nil {
		return responseErr.ResponseNotFound(TaskNotFound)
	}

	if err := instance.taskRepo.Delete(ctx, id); err != nil {
		instance.log.Error("failed to delete task by id ["+id+"] : ", zap.Error(err))
		return responseErr.ResponseInternalServerError(FailedToDeleteTask)
	}

	return nil
}

func (instance *taskService) GetAllWithPaginate(ctx context.Context, params *domain.TaskParams) (*domain.TaskPagination, error) {
	var (
		datas   []*domain.TaskTransformer
		maxPage int
	)

	if params.Limit > 100 {
		params.Limit = 100
	}

	tasks, total, err := instance.taskRepo.GetAllWithPaginate(ctx, params)
	if err != nil {
		instance.log.Error("failed to get task with limit : ", zap.Error(err))
		return nil, responseErr.ResponseInternalServerError(FailedToGetTask)
	}

	for _, task := range tasks {
		datas = append(datas, task.ToTaskTransformer())
	}

	if total%int64(params.Limit) > 0 {
		maxPage = int(total/int64(params.Limit)) + 1
	} else {
		maxPage = int(total / int64(params.Limit))
	}

	return &domain.TaskPagination{
		ListData: datas,
		PaginationData: domain.Pagination{
			CurrentPage:    params.Page,
			MaxDataPerPage: params.Limit,
			MaxPage:        maxPage,
			TotalAllData:   total,
		},
	}, nil
}
