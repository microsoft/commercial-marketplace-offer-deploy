package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"gorm.io/gorm"
)

type getDeploymentHandler struct {
	db     *gorm.DB
	mapper *mapper.DeploymentMapper
}

func (h *getDeploymentHandler) Handle(c echo.Context) error {
	id, err := h.getDeploymentId(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	deployment, err := h.getDeployment(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, h.mapper.Map(deployment))
}

func (h *getDeploymentHandler) getDeploymentId(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param(deploymenIdParameterName), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s invalid", deploymenIdParameterName)
	}
	return uint(id), nil
}

// method that gets a deployment struct by id
func (h *getDeploymentHandler) getDeployment(id uint) (*model.Deployment, error) {
	deployment := &model.Deployment{}
	result := h.db.First(deployment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return deployment, nil
}

// Get deployment factory function
func NewGetDeploymentHandler(appConfig *config.AppConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		handler := getDeploymentHandler{
			db:     db,
			mapper: mapper.NewDeploymentMapper(),
		}
		return handler.Handle(c)
	}
}
