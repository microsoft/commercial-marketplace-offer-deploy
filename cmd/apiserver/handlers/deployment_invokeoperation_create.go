package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

const deploymenIdParameterName = "deploymentId"

type invokeDeploymentOperation struct {
	repository operation.Repository
}

func (h *invokeDeploymentOperation) Handle(c echo.Context) error {
	deploymentId, request, err := h.getParameters(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	operation, err := h.createOperation(deploymentId, request)
	if err != nil {
		return err
	}

	err = operation.Schedule()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := sdk.InvokedDeploymentOperationResponse{
		InvokedOperation: &sdk.InvokedOperation{
			ID:         to.Ptr(operation.ID.String()),
			InvokedOn:  to.Ptr(time.Now().UTC()),
			Name:       request.Name,
			Parameters: request.Parameters,
			Result:     nil,
			Status:     to.Ptr(sdk.StatusScheduled.String()),
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *invokeDeploymentOperation) createOperation(deploymentId uint, request *sdk.InvokeDeploymentOperationRequest) (*operation.Operation, error) {
	operationType, configure, err := h.getConfigurator(deploymentId, request)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	operation, err := h.repository.New(operationType, configure)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return operation, nil
}

func (h *invokeDeploymentOperation) getConfigurator(deploymentId uint, request *sdk.InvokeDeploymentOperationRequest) (sdk.OperationType, operation.Configure, error) {
	if request == nil {
		return sdk.OperationUnknown, nil, fmt.Errorf("request is nil")
	}

	operationType, err := sdk.Type(*request.Name)
	if err != nil {
		return sdk.OperationUnknown, nil, err
	}

	parameters, ok := request.Parameters.(map[string]interface{})
	if !ok {
		parameters = make(map[string]interface{})
	}

	configure := func(i *model.InvokedOperation) error {
		retries := uint(*request.Retries)
		if retries <= 0 {
			retries = 1
		}
		i.Retries = retries
		i.DeploymentId = deploymentId
		i.Name = operationType.String()
		i.Parameters = parameters
		return nil
	}

	return operationType, configure, nil
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
		repository, err := operation.NewRepository(appConfig, nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		handler := invokeDeploymentOperation{
			repository: repository,
		}
		return handler.Handle(c)
	}
}
