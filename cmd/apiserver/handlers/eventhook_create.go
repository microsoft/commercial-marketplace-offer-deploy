package handlers

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo/v4"
	data "github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"gorm.io/gorm"
)

// HTTP handler for creating deployments
func CreateEventHook(c echo.Context, db *gorm.DB) error {
	var request *api.CreateEventHookRequest
	err := c.Bind(&request)

	if err != nil {
		return err
	}

	hook := &data.EventHook{}

	if hookExists(*request.Name, db) {
		db.Where(&data.EventHook{Name: *request.Name}).First(hook)
		hook.Callback = *request.Callback
		hook.ApiKey = *request.APIKey
		db.Save(&hook)

		result := &api.EventHookResponse{
			ID:       to.Ptr(hook.ID.String()),
			Name:     to.Ptr(hook.Name),
			Callback: to.Ptr(hook.Callback),
		}
		return c.JSON(http.StatusOK, result)
	}

	hook = data.FromCreateEventHook(request)
	tx := db.Save(&hook)
	if tx.Error != nil {
		return err
	}

	result := &api.EventHookResponse{
		ID:       to.Ptr(hook.ID.String()),
		Name:     &hook.Name,
		Callback: &hook.Callback,
	}

	return c.JSON(http.StatusOK, result)
}

func hookExists(hookName string, db *gorm.DB) bool {
	var count int64
	condition := data.EventHook{Name: hookName}
	db.Model(condition).Where(condition).Count(&count)
	return count > 0
}
