package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
)

type listDeploymentsHandler struct {
	query  *data.ListDeploymentsQuery
	mapper *mapper.DeploymentMapper
}

func (h *listDeploymentsHandler) Handle(c echo.Context) error {
	deployments := h.query.Execute()
	result := h.mapper.MapAll(deployments)
	return c.JSON(http.StatusOK, result)
}

func NewListDeploymentsHandler(appConfig *config.AppConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		handler := listDeploymentsHandler{
			query:  data.NewListDeploymentsQuery(db),
			mapper: mapper.NewDeploymentMapper(),
		}

		return handler.Handle(c)
	}
}
