package data

import (
	"errors"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"gorm.io/gorm"
)

type FindStageQuery struct {
	db *gorm.DB
}

// finds a stage by id. If found returns the parent deployment, the stage and nil error
// if not found, returns nil, nil, error
func (q *FindStageQuery) Execute(id uuid.UUID) (*model.Deployment, *model.Stage, error) {
	var list []model.Deployment
	q.db.Model(&model.Deployment{}).Find(&list)

	for _, deployment := range list {
		if len(deployment.Stages) > 0 {
			for _, stage := range deployment.Stages {
				if stage.ID == id {
					return &deployment, &stage, nil
				}
			}
		}
	}
	return nil, nil, errors.New("stage not found")
}

func NewFindStageQuery(db *gorm.DB) *FindStageQuery {
	return &FindStageQuery{
		db: db,
	}
}
