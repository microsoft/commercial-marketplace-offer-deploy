package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/persistence"
)

type HttpHandlerWithDatabase func(w http.ResponseWriter, r *http.Request, d persistence.Database)

// withDatabase wraps http handlers so a database is included as a func argument
func WithDatabase(handler HttpHandlerWithDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options := getDatabaseOptions()
		d := persistence.NewDatabase(options)
		handler(w, r, d)
	}
}

func getDatabaseOptions() *persistence.DatabaseOptions {
	configuration := config.GetConfiguration()
	dsn := filepath.Join(configuration.Database.Path, persistence.DatabaseFileName)
	options := &persistence.DatabaseOptions{Dsn: dsn, UseInMemory: configuration.Database.UseInMemory}

	return options
}
