package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/middleware"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(config *config.Configuration) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = addMiddleware(handler, route.Name, config)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

// default route
// TODO: update to provide a status / health check response
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Marketplace Offer Deployment Management Service\n-----------------------------------")

	for key, value := range r.Header {
		fmt.Fprintf(w, "\n"+key+" = "+strings.Join(value, ","))
	}
}

func addMiddleware(next http.Handler, routeName string, config *config.Configuration) http.Handler {
	handler := next
	//handler = middleware.AddJwtBearer(next, config)
	handler = middleware.AddLogging(handler, routeName)

	return handler
}
