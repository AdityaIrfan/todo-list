package ports

import (
	"context"
	"github.com/todo-list/internal/core/domain"
)

type (
	TaskService interface {
		Create(ctx context.Context, request *domain.CreateTaskRequst) error
		Update(ctx context.Context, id string, request *domain.UpdateTaskRequest) error
		Delete(ctx context.Context, id string) error
		GetOneByID(ctx context.Context, id string) (*domain.TaskTransformer, error)
		GetAllWithPaginate(ctx context.Context, params *domain.TaskParams) (*domain.TaskPagination, error)
	}
)
