package data

import "gorm.io/gorm"

// List deployments
type ListDeploymentsQuery struct {
	db *gorm.DB
}

func (q *ListDeploymentsQuery) Execute() []Deployment {
	var deployments []Deployment
	q.db.Find(&deployments)
	return deployments
}

type GetDeploymentQuery struct {
	db *gorm.DB
}

func (q *GetDeploymentQuery) Execute(id uint) Deployment {
	var deployment Deployment
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
