package data

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"gorm.io/gorm"
)

type StageQuery struct {
	db *gorm.DB
}

// finds a stage by id. If found returns the parent deployment, the stage and nil error
// if not found, returns nil, nil, error
func (q *StageQuery) Execute(id uuid.UUID) (*model.Deployment, *model.Stage, error) {
	stage := &model.Stage{}

	result := q.db.First(stage, id)

	if result.Error != nil || result.RowsAffected == 0 {
		err := result.Error
		if err == nil {
			err = fmt.Errorf("failed to get invoked operation: %v", result.Error)
		}
		return nil, nil, err
	}

	deployment := &model.Deployment{}
	q.db.First(deployment, stage.DeploymentID)
	return deployment, stage, nil
}

func NewStageQuery(db *gorm.DB) *StageQuery {
	return &StageQuery{
		db: db,
	}
}
