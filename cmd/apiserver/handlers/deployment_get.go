package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
