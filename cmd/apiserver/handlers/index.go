package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
)

type getHealthCheckHandler struct {
	mapper *mapper.GetHealthCheckResponseMapper
}

func (h *getHealthCheckHandler) Handle(c echo.Context) error {
	appInstance := hosting.GetApp()
	isHealthy := false
	if appInstance != nil {
		isHealthy = appInstance.IsReady()
	}
	return c.JSON(http.StatusOK, h.mapper.Map(isHealthy))
}

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Marketplace Offer Deployment Manager.")
}

func NewHealthCheckHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		handler := getHealthCheckHandler{
		}
		return handler.Handle(c)
	}
}
