package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func CreatEventSubscription(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func DeleteEventSubscription(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func GetEventSubscription(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func GetEvents(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func ListEventSubscriptions(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
