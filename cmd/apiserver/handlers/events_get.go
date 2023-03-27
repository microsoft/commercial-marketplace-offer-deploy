package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
)

// gets the event types
func GetEvents(c echo.Context) error {
	types := events.GetEventTypes()
	body := []generated.EventType{}

	for _, t := range types {
		body = append(body, generated.EventType{Name: &t})
	}
	return c.JSON(http.StatusOK, body)
}
