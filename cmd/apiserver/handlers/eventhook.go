package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func DeleteEventHook(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func GetEventHook(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
