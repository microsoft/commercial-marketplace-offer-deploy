package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"gorm.io/gorm"
)

type getDeploymentHandler struct {
	db     *gorm.DB
	mapper *mapper.DeploymentMapper
}

func (h *getDeploymentHandler) Handle(c echo.Context) error {
	//TODO: Get dep ID from echo param, query db for deployment (type data.deployment)
	//Get Deployment ID from echo param
	deploymentId, err := strconv.Atoi(c.Param(deploymenIdParameterName))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s invalid", deploymenIdParameterName))
	}

	//get deployment by id
	data := &data.Deployment{}
	deployment := h.db.First(data, deploymentId)
	if deployment.Error != nil {
		return err
	}

	//map data.deployment to api.deployment
	h.mapper.Map(deployment)
	result := createResult(deployment)
	//return JSON as api.deployment type
	return c.JSON(http.StatusOK, result)
}

// Get deployment factory function
func NewGetDeploymentHandler(appConfig *config.AppConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		//construct gorm.db instance using appconfig
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		handler := getDeploymentHandler{
			db:     db,
			mapper: mapper.NewDeploymentMapper(),
		}
		return handler.Handle(c)
	}
}
