package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
