package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"gorm.io/gorm"
)

type listInvokedOperationHandler struct {
	db     *gorm.DB
	mapper *mapper.InvokedDeploymentResponseMapper
}

func (h *listInvokedOperationHandler) Handle(c echo.Context) error {
	list, err := h.list()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, h.mapper.MapList(list))
}

// method that gets a deployment struct by id
func (h *listInvokedOperationHandler) list() ([]model.InvokedOperation, error) {
	list := []model.InvokedOperation{}
	h.db.Find(&list)
	return list, h.db.Error
}

func NewListInvokedOperationHandler(appConfig *config.AppConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		handler := listInvokedOperationHandler{
			db:     db,
			mapper: &mapper.InvokedDeploymentResponseMapper{},
		}
		return handler.Handle(c)
	}
}
