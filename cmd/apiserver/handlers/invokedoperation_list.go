package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
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

func (h *listInvokedOperationHandler) getId(c echo.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(operationIdParameterName))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s invalid", operationIdParameterName)
	}
	return id, nil
}

// method that gets a deployment struct by id
func (h *listInvokedOperationHandler) list() ([]data.InvokedOperation, error) {
	list := []data.InvokedOperation{}
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
