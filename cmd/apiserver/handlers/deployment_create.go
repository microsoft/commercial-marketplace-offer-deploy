package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"gorm.io/gorm"
)

type createDeploymentHandler struct {
	db     *gorm.DB
	mapper *mapper.CreateDeploymentMapper
}

// HTTP handler for creating deployments
func (h *createDeploymentHandler) Handle(c echo.Context) error {
	var command *api.CreateDeployment
	err := c.Bind(&command)

	if err != nil {
		return err
	}

	deployment, err := h.mapper.Map(command)
	if err != nil {
		return err
	}

	tx := h.db.Create(&deployment)
	log.Debug("Deployment [%d] created.", deployment.ID)

	if tx.Error != nil {
		return err
	}

	if err != nil {
		return err
	}

	result := createResult(deployment)
	return c.JSON(http.StatusOK, result)
}

func createResult(deployment *data.Deployment) *api.Deployment {
	result := &api.Deployment{
		ID:     to.Ptr(int32(deployment.ID)),
		Name:   &deployment.Name,
	}

	for _, stage := range deployment.Stages {
		result.Stages = append(result.Stages, &api.DeploymentStage{
			Name: to.Ptr(stage.Name),
			ID:   to.Ptr(stage.ID.String()),
			Status: &stage.Status,
			DeploymentName: &stage.DeploymentName,
		})
	}
	return result
}

func NewCreateDeploymentHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()

		handler := createDeploymentHandler{
			db:     db,
			mapper: mapper.NewCreateDeploymentMapper(),
		}
		return handler.Handle(c)
	}
}
