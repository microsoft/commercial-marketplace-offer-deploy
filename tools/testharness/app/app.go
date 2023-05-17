package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
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
)

func AddRoutes(e *echo.Echo) {
	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Test Harness Up.")
	})
	e.GET("/status/:deploymentId/:operationName", GetStatus)
	e.GET("/setcase/:caseName", SetCase)
	e.GET("/createdeployment", CreateDeployment)
	e.GET("/startdeployment/:deploymentId", StartDeployment)
	e.GET("/createeventhook", CreateEventHook)
	e.GET("/dryrun/:deploymentId", DryRun)
	e.GET("/redeploy/:deploymentId/:stageName", Redeploy)
	e.POST("/webhook", ReceiveEventHook)
}
