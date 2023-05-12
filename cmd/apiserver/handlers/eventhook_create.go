package handlers

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	data "github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"gorm.io/gorm"
)

type createEventHookHandler struct {
	db *gorm.DB
}

// HTTP handler for creating deployments
func (h createEventHookHandler) Handle(c echo.Context) error {
	db := h.db
	var request *api.CreateEventHookRequest
	err := c.Bind(&request)

	if err != nil {
		return err
	}

	hook := &data.EventHook{}

	if h.hookExists(*request.Name, db) {
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

func (h *createEventHookHandler) hookExists(hookName string, db *gorm.DB) bool {
	var count int64
	condition := data.EventHook{Name: hookName}
	db.Model(condition).Where(condition).Count(&count)
	return count > 0
}

// NewCreateEventHookHandler creates a new instance of the createEventHookHandler
func NewCreateEventHookHandler(appConfig *config.AppConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		d := data.NewDatabase(appConfig.GetDatabaseOptions())
		h := createEventHookHandler{db: d.Instance()}
		return h.Handle(c)
	}
}
