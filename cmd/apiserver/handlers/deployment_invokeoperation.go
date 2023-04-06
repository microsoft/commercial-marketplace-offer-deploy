package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	data "github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	ops "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
	"gorm.io/gorm"
)

const deploymenIdParameterName = "deploymentId"

type invokeDeploymentOperation struct {
	db              *gorm.DB
	operationHandle operations.DeploymentOperationHandlerFunc
}

func (h *invokeDeploymentOperation) Handle(c echo.Context) error {
	deploymentId, err := strconv.Atoi(c.Param(deploymenIdParameterName))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s invalid", deploymenIdParameterName))
	}

	var operation api.InvokeDeploymentOperation
	err = c.Bind(&operation)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	name, err := ops.From(*operation.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	context := h.getContext(deploymentId, name, operation, c)
	result, err := h.operationHandle(context)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (*invokeDeploymentOperation) getContext(deploymentId int, name ops.OperationType, operation api.InvokeDeploymentOperation, c echo.Context) *operations.DeploymentOperationContext {
	context := &operations.DeploymentOperationContext{
		DeploymentId: deploymentId,
		Name:         name,
		Parameters:   operation.Parameters.(operations.DeploymentOperationParameters),
		HttpContext:  c,
	}
	return context
}

func NewInvokeDeploymentOperationHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) echo.HandlerFunc {
	return func(c echo.Context) error {
		handler := invokeDeploymentOperation{
			db:              data.NewDatabase(appConfig.GetDatabaseOptions()).Instance(),
			operationHandle: operations.NewStartDeploymentOperationHandler(appConfig, credential),
		}
		return handler.Handle(c)
	}
}
