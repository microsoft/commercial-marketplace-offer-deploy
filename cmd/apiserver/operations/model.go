package operations

import "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"

type InvokeOperationCommand struct {
	DeploymentId int
	Request      *api.InvokeDeploymentOperationRequest
}
