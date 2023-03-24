package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
	"gorm.io/gorm"
)

const DeploymenIdParameterName = "deploymentId"

type InvokeOperationDeploymentHandler func(int, generated.InvokeDeploymentOperation, *gorm.DB) (interface{}, error)

func GetDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func InvokeOperation(c echo.Context, db *gorm.DB) error {
	deploymentId, err := strconv.Atoi(c.Param(DeploymenIdParameterName))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s in route was not an int", DeploymenIdParameterName))
	}

	var operation generated.InvokeDeploymentOperation
	err = c.Bind(&operation)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	log.Printf("Operation deserialized \n %v", operation)

	operationHandler := CreateOperationHandler(operation)

	if operationHandler == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "There was op OperationHandler registered for the invoked operation")
	}
	res, err := operationHandler(deploymentId, operation, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func CreateOperationHandler(operation generated.InvokeDeploymentOperation) InvokeOperationDeploymentHandler {
	switch operation.Name {
	case operation.Name:
		return CreateDryRun
	default:
		return nil
	}
}

func ListDeployments(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func UpdateDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
