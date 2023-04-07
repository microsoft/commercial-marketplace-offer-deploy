package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	ops "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

const deploymenIdParameterName = "deploymentId"

type invokeDeploymentOperation struct {
	processor operations.InvokeOperationProcessor
}

func (h *invokeDeploymentOperation) Handle(c echo.Context) error {
	deploymentId, err := strconv.Atoi(c.Param(deploymenIdParameterName))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s invalid", deploymenIdParameterName))
	}

	var request *api.InvokeDeploymentOperationRequest
	err = c.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	command := operations.InvokeOperationCommand{
		DeploymentId: deploymentId,
		Request:      request,
	}
	operationId, err := h.processor.Process(c.Request().Context(), &command)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := &api.InvokedDeploymentOperationResponse{
		ID:         to.Ptr(operationId.String()),
		InvokedOn:  to.Ptr(time.Now().UTC()),
		Name:       request.Name,
		Parameters: request.Parameters,
		Result:     ops.OperationResultAccepted,
		Status:     to.Ptr(events.DeploymentPendingEventType.String()),
	}

	return c.JSON(http.StatusOK, response)
}

func NewInvokeDeploymentOperationHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) echo.HandlerFunc {
	return func(c echo.Context) error {
		processor, err := operations.NewInvokeOperationProcessor(appConfig, credential)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		handler := invokeDeploymentOperation{
			processor: processor,
		}
		return handler.Handle(c)
	}
}
