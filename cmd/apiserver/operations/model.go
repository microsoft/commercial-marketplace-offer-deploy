package operations

import (
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

type DeploymentOperationParameters map[string]any

type DeploymentOperationContext struct {
	Name         operations.OperationType
	DeploymentId int
	Parameters   DeploymentOperationParameters
	HttpContext  echo.Context
}

// deployment operation handler func
type DeploymentOperationHandlerFunc func(c *DeploymentOperationContext) (*api.InvokedDeploymentOperation, error)
