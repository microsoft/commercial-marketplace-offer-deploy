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

type listEventHooksHandler struct {
	db *gorm.DB
}

func (h *listEventHooksHandler) Handle(c echo.Context) error {
	hooks := []data.EventHook{}
	h.db.Find(&hooks)

	result := []api.EventHookResponse{}
	for _, hook := range hooks {
		result = append(result, api.EventHookResponse{
			ID:       to.Ptr(hook.ID.String()),
			Name:     to.Ptr(hook.Name),
			Callback: to.Ptr(hook.Callback),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func NewListEventHooksHandler(appConfig *config.AppConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		handler := listEventHooksHandler{
			db: data.NewDatabase(appConfig.GetDatabaseOptions()).Instance(),
		}
		return handler.Handle(c)
	}

}
