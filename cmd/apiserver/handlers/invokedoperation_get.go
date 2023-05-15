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

const operationIdParameterName = "operationId"

type getInvokedOperationHandler struct {
	db     *gorm.DB
	mapper *mapper.InvokedDeploymentOperationResponseMapper
}

func (h *getInvokedOperationHandler) Handle(c echo.Context) error {
	id, err := h.getId(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	invokedOperation, err := h.get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, h.mapper.Map(invokedOperation))
}

func (h *getInvokedOperationHandler) getId(c echo.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(operationIdParameterName))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s invalid", operationIdParameterName)
	}
	return id, nil
}

// method that gets a deployment struct by id
func (h *getInvokedOperationHandler) get(id uuid.UUID) (*data.InvokedOperation, error) {
	invokedOperation := &data.InvokedOperation{}
	result := h.db.First(invokedOperation, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return invokedOperation, nil
}

func NewGetInvokedOperationHandler(appConfig *config.AppConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		handler := getInvokedOperationHandler{
			db:     db,
			mapper: &mapper.InvokedDeploymentOperationResponseMapper{},
		}
		return handler.Handle(c)
	}
}
