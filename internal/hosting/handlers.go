package hosting

import (
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"gorm.io/gorm"
)

// Defines an HTTP handler that also has a db instance as a parameter
type DataHandlerFunc func(c echo.Context, db *gorm.DB) error

// Wraps a data handler func into an echo.HandlerFunc for route registration purposes
func ToHandlerFunc(h DataHandlerFunc, databaseOptions *data.DatabaseOptions) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Printf("Inside ToHandlerFunc with database options %v", databaseOptions)
		d := data.NewDatabase(databaseOptions)
		return h(c, d.Instance())
	}
}
