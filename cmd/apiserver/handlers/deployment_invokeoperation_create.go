package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/dispatch"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
)

const deploymenIdParameterName = "deploymentId"

type invokeDeploymentOperation struct {
	dispatcher dispatch.OperatorDispatcher
}

func (h *invokeDeploymentOperation) Handle(c echo.Context) error {
	deploymentId, request, err := h.getParameters(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	operationId, err := h.dispatcher.Dispatch(c.Request().Context(), &dispatch.DispatchInvokedOperation{
		DeploymentId: deploymentId,
		Request:      request,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := api.InvokedDeploymentOperationResponse{
		InvokedOperation: &api.InvokedOperation{
			ID:         to.Ptr(operationId.String()),
			InvokedOn:  to.Ptr(time.Now().UTC()),
			Name:       request.Name,
			Parameters: request.Parameters,
			Result:     nil,
			Status:     to.Ptr(operation.StatusScheduled.String()),
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *invokeDeploymentOperation) getParameters(c echo.Context) (int, *api.InvokeDeploymentOperationRequest, error) {
	deploymentId, err := strconv.Atoi(c.Param(deploymenIdParameterName))
	if err != nil {
		return deploymentId, nil, fmt.Errorf("%s invalid", deploymenIdParameterName)
	}

	var request *api.InvokeDeploymentOperationRequest
	err = c.Bind(&request)
	if err != nil {
		return deploymentId, nil, err
	}

	return deploymentId, request, nil
}

func NewInvokeDeploymentOperationHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) echo.HandlerFunc {
	return func(c echo.Context) error {
		processor, err := dispatch.NewOperatorDispatcher(appConfig, credential)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		handler := invokeDeploymentOperation{
			dispatcher: processor,
		}
		return handler.Handle(c)
	}
}
