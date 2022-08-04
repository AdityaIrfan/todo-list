package taskrps

import (
	"context"
	"errors"
	"github.com/todo-list/internal/core/domain"
	"github.com/todo-list/internal/core/ports"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type taskPostgres struct {
	postgres *gorm.DB
}

func NewTaskPostgres(postgres *gorm.DB) ports.TaskRepository {
	return &taskPostgres{
		postgres: postgres,
	}
}

// Create is creating new task and objectives
func (instance *taskPostgres) Create(ctx context.Context, task *domain.Task) error {
	if err := instance.postgres.Debug().Transaction(func(tx *gorm.DB) error {
		// save task
		if err := tx.Debug().Save(&task).Error; err != nil {
			return err
		}

		// save objectives
		for _, obj := range task.Objective {
			obj.TaskID = task.ID
			if err := tx.Debug().Save(&obj).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Update is updating task and objectives
func (instance *taskPostgres) Update(ctx context.Context, task *domain.Task) error {
	if err := instance.postgres.Debug().Transaction(func(tx *gorm.DB) error {
		if len(task.Objective) > 0 {
			// delete objectives
			if err := tx.Debug().
				Where("task_id = ?", strconv.Itoa(int(task.ID))).
				Delete(&domain.Objective{}).Unscoped().Error; err != nil {
				return err
			}

			// save new objectives
			for _, obj := range task.Objective {
				obj.TaskID = task.ID
				if err := tx.Debug().Save(&obj).Error; err != nil {
					return err
				}
			}
		}

		// update task
		if err := tx.Save(&task).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetOneByID is getting task by id and its objectives
func (instance *taskPostgres) GetOneByID(ctx context.Context, id string) (*domain.Task, error) {
	var task *domain.Task

	if err := instance.postgres.Debug().Preload("Objective").Where("id = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return task, nil
}

func (instance *taskPostgres) Delete(ctx context.Context, id string) error {
	if err := instance.postgres.Transaction(func(tx *gorm.DB) error {
		// delete objective
		if err := tx.Where("task_id = ?", id).Delete(&domain.Objective{}).Unscoped().Error; err != nil {
			return err
		}

		// delete task
		if err := tx.Where("id = ?", id).Delete(&domain.Task{}).Unscoped().Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (instance *taskPostgres) GetAllWithPaginate(ctx context.Context, params *domain.TaskParams) ([]*domain.Task, int64, error) {
	var (
		tasks []*domain.Task
		total int64
	)

	q := instance.postgres.Preload("Objective").Debug()

	if params.Title != nil {
		q = q.Where(`LOWER(title) LIKE LOWER(?)`, "%"+*params.Title+"%")
	}
	if params.ActionTimeStart != nil {
		q = q.Where("action_time >= ?", time.Unix(int64(*params.ActionTimeStart), 0).UTC())
	}
	if params.ActionTimeEnd != nil {
		q = q.Where("action_time <= ?", time.Unix(int64(*params.ActionTimeEnd), 0).UTC())
	}
	if params.IsFinished != nil {
		q = q.Where("is_finished = ?", params.IsFinished)
	}

	if err := q.Model(&domain.Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Limit(params.Limit).Offset(params.Limit * (params.Page - 1)).Find(&tasks).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	return tasks, total, nil
}
