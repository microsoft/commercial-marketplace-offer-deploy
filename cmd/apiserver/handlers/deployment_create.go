package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	data "github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"gorm.io/gorm"
)

// HTTP handler for creating deployments
func CreateDeployment(c echo.Context, db *gorm.DB) error {
	log.Println("Inside createdeplyoment")
	var command *api.CreateDeployment
	err := c.Bind(&command)

	if err != nil {
		return err
	}

	deployment := data.FromCreateDeployment(command)
	log.Printf("Deployment: %v", deployment)
	tx := db.Create(&deployment)

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
