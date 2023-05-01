package sdk

import (
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
)

type RetryOptions struct {
	StageId uuid.UUID
}

type StatusResponse struct {
	Id           uuid.UUID
	Name         string
	Status       string
	Result       any
	Attempts     int
	DeploymentId int
}

type DryRunResponse struct {
	Id uuid.UUID
	//Results map[string]any
	Results any
	Status  string
}

type StartOptions struct {
	Retries int
}

type StartDeploymentResponse struct {
	Id     uuid.UUID
	Status string
}

type RetryResponse struct {
	Id         uuid.UUID
	Status     string
	Parameters map[string]any
}

type ListResponse struct {
	Deployments []*api.Deployment
}
