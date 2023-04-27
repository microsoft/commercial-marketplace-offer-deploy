package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Marketplace Offer Deployment Manager.")
}
