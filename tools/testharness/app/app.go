package app

import (
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/tools/testharness/audit"
)

// TODO: this needs to go and pull from .env
var (
	location       = "eastus"
	caseName       = "success"
	resourceGroup  string
	subscription   string
	clientEndpoint = "http://localhost:8080"
	env            = loadEnvironmentVariables()
	deployment     *sdk.Deployment
	messageAudit   = audit.NewAppendOnlyFileAuditLog("/tmp/testharness.eventhookmessages.txt")
)

func AddRoutes(e *echo.Echo) {
	e.GET("/", HealthStatus)
	e.GET("/status/:deploymentId/:operationName", GetStatus)
	e.GET("/setcase/:caseName", SetCase)
	e.GET("/createdeployment", CreateDeployment)
	e.GET("/getdeployment/:deploymentId", GetDeployment)
	e.GET("/startdeployment/:deploymentId", StartDeployment)
	e.GET("/createeventhook", CreateEventHook)
	e.GET("/dryrun/:deploymentId", DryRun)
	e.GET("/redeploy/:deploymentId/:stageName", Redeploy)
	e.GET("/cancel/:deploymentId", Cancel)
	e.POST("/webhook", ReceiveEventHook)
}
