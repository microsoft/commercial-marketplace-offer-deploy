package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
