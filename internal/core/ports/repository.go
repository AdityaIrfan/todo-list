package ports

import (
	"context"
	"github.com/todo-list/internal/core/domain"
)

type (
	TaskRepository interface {
		Create(ctx context.Context, task *domain.Task) error
		Update(ctx context.Context, task *domain.Task) error
		Delete(ctx context.Context, id string) error
		GetOneByID(ctx context.Context, id string) (*domain.Task, error)
		GetAllWithPaginate(ctx context.Context, params *domain.TaskParams) ([]*domain.Task, int64, error)
	}
)
