package data

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"gorm.io/gorm"
)

// List deployments
type ListDeploymentsQuery struct {
	db *gorm.DB
}

func (q *ListDeploymentsQuery) Execute() []model.Deployment {
	var deployments []model.Deployment
	q.db.Find(&deployments)
	return deployments
}

type GetDeploymentQuery struct {
	db *gorm.DB
}

func (q *GetDeploymentQuery) Execute(id uint) model.Deployment {
	var deployment model.Deployment
	q.db.First(&deployment)
	return deployment
}

//region factory

func NewListDeploymentsQuery(db *gorm.DB) *ListDeploymentsQuery {
	return &ListDeploymentsQuery{
		db: db,
	}
}

func NewGetDeploymentQuery(db *gorm.DB) *GetDeploymentQuery {
	return &GetDeploymentQuery{
		db: db,
	}
}

//endregion factory
