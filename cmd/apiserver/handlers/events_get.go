package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
)

// gets the event types
func GetEvents(c echo.Context) error {
	types := events.GetEventTypes()
	body := []api.EventType{}

	for _, t := range types {
		body = append(body, api.EventType{Name: &t})
	}
	return c.JSON(http.StatusOK, body)
}
