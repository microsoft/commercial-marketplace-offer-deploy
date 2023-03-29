package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/azure/eventgrid"
	"gorm.io/gorm"
)

// this handler is the webook endpoint for event grid events
func EventGridWebHook(c echo.Context, db *gorm.DB) error {
	webhookValidator := eventgrid.NewWebHookValidationEventHandler(c.Bind)
	result := webhookValidator.Handle(c.Request())

	if result.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
	}

	if result.Handled {
		return c.JSON(http.StatusOK, &result.Response)
	}

	//TODO: handle all other event types
	return nil
}
