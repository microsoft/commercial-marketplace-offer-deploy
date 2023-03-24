package handlers

import (
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/runtime"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"gorm.io/gorm"
)

type DataHandlerFunc func(c echo.Context, db *gorm.DB) error

func ToHandlerFunc(h DataHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		options := createOptionsFromConfiguration()
		d := data.NewDatabase(options)
		return h(c, d.Instance())
	}
}

func createOptionsFromConfiguration() *data.DatabaseOptions {
	configuration := runtime.GetApp().GetConfiguration()
	dsn := filepath.Join(configuration.Database.Path, data.DatabaseFileName)
	options := &data.DatabaseOptions{Dsn: dsn, UseInMemory: configuration.Database.UseInMemory}

	return options
}
