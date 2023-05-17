package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// gets the event types
func GetEvents(c echo.Context) error {
	types := sdk.GetEventTypes()
	body := []sdk.EventType{}

	for _, t := range types {
		body = append(body, sdk.EventType{Name: &t})
	}
	return c.JSON(http.StatusOK, body)
}
