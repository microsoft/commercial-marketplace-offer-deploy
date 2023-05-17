package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/dispatch"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	"gorm.io/gorm"
)

const deploymenIdParameterName = "deploymentId"

type invokeDeploymentOperation struct {
	db     *gorm.DB
	dispatcher dispatch.OperatorDispatcher
}

func (h *invokeDeploymentOperation) Handle(c echo.Context) error {
	deploymentId, request, err := h.getParameters(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	isValid, err := h.validateParameters(deploymentId, request.Parameters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if !isValid {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("invalid parameters"))
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

func (h *invokeDeploymentOperation) validateParameters(deploymentId int, parameters any) (bool, error) {
	deployment, err := h.getDeployment(deploymentId)
	if err != nil {
		return false, err
	}
	if mapParams, ok := parameters.(map[string]any); ok {
		if deployment.Template != nil && deployment.Template["parameters"] != nil {
			templateParameters := deployment.Template["parameters"].(map[string]any)
			hasAllKeys := utils.CheckKeys(templateParameters, mapParams)
			return hasAllKeys, nil
		}
	}
	return false, errors.New("could not propertly validateParamsters of invoked operation against deployment")
}

func (h *invokeDeploymentOperation) getDeployment(id int) (*data.Deployment, error) {
	deployment := &data.Deployment{}
	result := h.db.First(deployment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return deployment, nil
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

		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		handler := invokeDeploymentOperation{
			dispatcher: processor,
			db:         db,
		}
		return handler.Handle(c)
	}
}
