package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Marketplace Offer Deployment Manager.")
}
