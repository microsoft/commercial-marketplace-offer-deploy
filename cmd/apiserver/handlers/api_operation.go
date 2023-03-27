package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetDeploymentOperation(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func ListOperations(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
