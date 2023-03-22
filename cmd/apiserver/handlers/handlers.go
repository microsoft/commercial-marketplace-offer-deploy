package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
)

type HttpHandlerWithDatabase func(w http.ResponseWriter, r *http.Request, d *data.Database)

// withDatabase wraps http handlers so a database is included as a func argument
func WithDatabase(handler HttpHandlerWithDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options := createOptionsFromConfiguration()
		d := data.NewDatabase(options)
		handler(w, r, &d)
	}
}

func createOptionsFromConfiguration() *data.DatabaseOptions {
	configuration := config.GetConfiguration()
	dsn := filepath.Join(configuration.Database.Path, data.DatabaseFileName)
	options := &data.DatabaseOptions{Dsn: dsn, UseInMemory: configuration.Database.UseInMemory}

	return options
}
