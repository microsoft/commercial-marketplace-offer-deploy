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
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
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

	command, err := h.createCommand(deploymentId, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	operationId, err := h.dispatcher.Dispatch(c.Request().Context(), command)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := sdk.InvokedDeploymentOperationResponse{
		InvokedOperation: &sdk.InvokedOperation{
			ID:         to.Ptr(operationId.String()),
			InvokedOn:  to.Ptr(time.Now().UTC()),
			Name:       request.Name,
			Parameters: request.Parameters,
			Result:     nil,
			Status:     to.Ptr(sdk.StatusScheduled.String()),
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *invokeDeploymentOperation) createCommand(deploymentId uint, request *sdk.InvokeDeploymentOperationRequest) (*dispatch.DispatchInvokedOperation, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	parameters, ok := request.Parameters.(map[string]interface{})
	if !ok {
		parameters = make(map[string]interface{})
	}

	command := &dispatch.DispatchInvokedOperation{
		DeploymentId: deploymentId,
		Retries:      int(*request.Retries),
		Name:         *request.Name,
		Parameters:   parameters,
	}

	if command.Retries <= 0 {
		command.Retries = 1
	}

	return command, nil
}

func (h *invokeDeploymentOperation) getParameters(c echo.Context) (uint, *sdk.InvokeDeploymentOperationRequest, error) {
	deploymentId, err := strconv.Atoi(c.Param(deploymenIdParameterName))
	if err != nil {
		return uint(deploymentId), nil, fmt.Errorf("%s invalid", deploymenIdParameterName)
	}

	var request *sdk.InvokeDeploymentOperationRequest
	err = c.Bind(&request)
	if err != nil {
		return uint(deploymentId), nil, err
	}

	return uint(deploymentId), request, nil
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
