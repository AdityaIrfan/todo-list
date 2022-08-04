package domain

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Task struct {
	ID         uint64
	Title      string
	ActionTime time.Time
	IsFinished bool
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Objective []*Objective
}

func (t *Task) GetObjectives() []ObjectiveTransformer {
	var transformer []ObjectiveTransformer

	for _, objective := range t.Objective {
		transformer = append(transformer, *objective.ToObjectiveTransformer())
	}

	return transformer
}

type TaskTransformer struct {
	ID         uint64                 `json:"Task_ID"`
	Title      string                 `json:"Title"`
	ActionTime int64                  `json:"Action_Time"`
	CreatedAt  int64                  `json:"Created_Time"`
	UpdatedAt  int64                  `json:"Updated_Time"`
	IsFinished bool                   `json:"Is_Finished"`
	Objectives []ObjectiveTransformer `json:"Objective_List"`
}

func (t *Task) ToTaskTransformer() *TaskTransformer {
	return &TaskTransformer{
		ID:         t.ID,
		Title:      t.Title,
		ActionTime: t.ActionTime.Unix(),
		CreatedAt:  t.CreatedAt.Unix(),
		UpdatedAt:  t.UpdatedAt.Unix(),
		IsFinished: t.IsFinished,
		Objectives: t.GetObjectives(),
	}
}

type CreateTaskRequst struct {
	Title      string   `json:"Title"`
	ActionTime int64    `json:"Action_Time"`
	Objectives []string `json:"Objective_List"`
}

func (c *CreateTaskRequst) ToBaseObjectives() []*Objective {
	var objectives []*Objective

	for _, obj := range c.Objectives {
		objectives = append(objectives, &Objective{
			ObjectiveName: obj,
			IsFinished:    false,
		})
	}

	return objectives
}

func (c *CreateTaskRequst) ToBase() *Task {
	return &Task{
		Title:      c.Title,
		ActionTime: time.Unix(c.ActionTime, 0).UTC(),
		IsFinished: false,
		Objective:  c.ToBaseObjectives(),
	}
}

type UpdateTaskRequest struct {
	Title      string                   `json:"title"`
	Objectives []UpdateObjectiveRequest `json:"Objective_List"`
}

func (u *UpdateTaskRequest) ToBaseObjectives() ([]*Objective, bool) {
	var (
		objectives    []*Objective
		isAllFinished bool = false
		finishedCount int  = 0
	)

	for _, obj := range u.Objectives {
		if obj.IsFinished {
			finishedCount++
		}

		objectives = append(objectives, &Objective{
			ObjectiveName: obj.ObjectiveName,
			IsFinished:    obj.IsFinished,
		})
	}

	if finishedCount == len(u.Objectives) {
		isAllFinished = true
	}

	return objectives, isAllFinished
}

func (u *UpdateTaskRequest) ToBase(task *Task) *Task {
	objestives, isAllFinished := u.ToBaseObjectives()

	return &Task{
		ID:         task.ID,
		Title:      u.Title,
		ActionTime: task.ActionTime,
		IsFinished: isAllFinished,
		CreatedAt:  task.CreatedAt,
		Objective:  objestives,
	}
}

type TaskParams struct {
	Page            int     `query:"Page"`
	Limit           int     `query:"Limit"`
	Title           *string `query:"Title"`
	ActionTimeStart *int    `query:"Action_Time_Start"`
	ActionTimeEnd   *int    `query:"Action_Time_End"`
	IsFinished      *bool   `query:"Is_Finished"`
}

func (t TaskParams) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Page, validation.Required, validation.By(moreThanNol)),
		validation.Field(&t.Limit, validation.Required, validation.By(moreThanNol)),
	)
}

func moreThanNol(value interface{}) error {
	intValue, _ := value.(int)

	if intValue < 0 {
		return errors.New("value must be more than 0")
	}

	return nil
}

type TaskPagination struct {
	ListData       []*TaskTransformer `json:"List_Data"`
	PaginationData Pagination         `json:"Pagination_Data"`
}
