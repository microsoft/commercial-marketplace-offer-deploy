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
func CreateEventSubscription(c echo.Context, db *gorm.DB) error {
	eventType := c.Param("eventType")

	var request *api.CreateEventSubscriptionRequest
	err := c.Bind(&request)

	if err != nil {
		return err
	}

	model := data.FromCreateEventSubscription(eventType, request)
	tx := db.Create(&model)

	log.Printf("Event Subscription [%s] created for event type %s.", model.Name, model.EventType)

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
