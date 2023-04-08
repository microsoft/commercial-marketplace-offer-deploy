package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func ListDeployments(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
