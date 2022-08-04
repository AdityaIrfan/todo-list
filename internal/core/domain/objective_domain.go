package domain

type Objective struct {
	ID            uint64
	TaskID        uint64
	ObjectiveName string
	IsFinished    bool

	Task *Task `gorm:"foreignKey:TaskID;references:ID"`
}

type ObjectiveTransformer struct {
	ObjectiveName string `json:"Objective_Name"`
	IsFinished    bool   `json:"Is_Finished"`
}

func (o *Objective) ToObjectiveTransformer() *ObjectiveTransformer {
	return &ObjectiveTransformer{
		ObjectiveName: o.ObjectiveName,
		IsFinished:    o.IsFinished,
	}
}

type UpdateObjectiveRequest struct {
	ObjectiveName string `json:"Objective_Name"`
	IsFinished    bool   `json:"Is_Finished"`
}
