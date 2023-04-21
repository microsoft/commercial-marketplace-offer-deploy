package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func Index(c echo.Context) error {
	logrus.Info("Test from the Index")
	return c.String(http.StatusOK, "Marketplace Offer Deployment Management Service\n-----------------------------------")
}
