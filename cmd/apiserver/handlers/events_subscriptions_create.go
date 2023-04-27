package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	data "github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"gorm.io/gorm"
)

// HTTP handler for creating deployments
func CreateEventSubscription(c echo.Context, db *gorm.DB) error {
	var request *api.CreateEventSubscriptionRequest
	err := c.Bind(&request)

	if err != nil {
		return err
	}

	// TODO: validate with a test handshake before continuing

	model := data.FromCreateEventSubscription(request)
	tx := db.Create(&model)

	if tx.Error != nil {
		return err
	}

	id := model.ID.String()
	result := &api.EventSubscriptionResponse{
		ID:       &id,
		Name:     &model.Name,
		Callback: &model.Callback,
	}

	return c.JSON(http.StatusOK, result)
}
