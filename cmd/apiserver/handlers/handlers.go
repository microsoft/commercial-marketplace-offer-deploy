package handlers

import (
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"gorm.io/gorm"
)

type DataHandlerFunc func(c echo.Context, db *gorm.DB) error

func ToHandlerFunc(h DataHandlerFunc, configuration *config.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		options := createOptionsFromConfiguration(configuration)
		d := data.NewDatabase(options)
		return h(c, d.Instance())
	}
}

func createOptionsFromConfiguration(c *config.Configuration) *data.DatabaseOptions {
	dsn := filepath.Join(c.Database.Path, data.DatabaseFileName)
	options := &data.DatabaseOptions{Dsn: dsn, UseInMemory: c.Database.UseInMemory}

	return options
}
