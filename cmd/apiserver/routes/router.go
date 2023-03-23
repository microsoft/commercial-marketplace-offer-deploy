package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

type Route struct {
	Method      string
	Name        string
	Path        string
	HandlerFunc echo.HandlerFunc
}

type Routes []Route

// func NewRouter(config *config.Configuration) *mux.Router {
// 	router := mux.NewRouter().StrictSlash(true)
// 	for _, route := range routes {
// 		var handler http.Handler
// 		handler = route.HandlerFunc
// 		handler = addMiddleware(handler, route.Name, config)

// 		router.
// 			Methods(route.Method).
// 			Path(route.Path).
// 			Name(route.Name).
// 			Handler(handler)
// 	}

// 	return router
// }

// default route
// TODO: update to provide a status / health check response
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Marketplace Offer Deployment Management Service\n-----------------------------------")
}

// func addMiddleware(next http.Handler, routeName string, config *config.Configuration) http.Handler {
// 	handler := next
// 	//handler = middleware.AddJwtBearer(next, config)
// 	handler = middleware.AddLogging(handler, routeName)

// 	return handler
// }
