package handlers

import (
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/labstack/echo"
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
	log.Printf("Deployment [%d] created.", deployment.ID)

	if tx.Error != nil {
		return err
	}

	if err != nil {
		return err
	}

	deploymentId := int32(deployment.ID)
	result := api.Deployment{
		ID:     &deploymentId,
		Name:   &deployment.Name,
		Status: &deployment.Status,
	}

	return c.JSON(http.StatusOK, result)
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
