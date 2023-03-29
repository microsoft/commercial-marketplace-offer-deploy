package handlers

import (
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// this handler is the webook endpoint for event grid events
func EventGridWebHook(c echo.Context, db *gorm.DB) error {
	return nil
}
