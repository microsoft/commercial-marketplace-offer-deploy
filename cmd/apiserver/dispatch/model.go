package dispatch

import "github.com/microsoft/commercial-marketplace-offer-deploy/sdk"

// We want to take a request to invoke and operation and dispatch.
// This is the command to do that.
type DispatchInvokedOperation struct {
	DeploymentId int
	Request      *sdk.InvokeDeploymentOperationRequest
}
