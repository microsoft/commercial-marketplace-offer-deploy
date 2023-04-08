package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func UpdateDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
