package hosting

import "github.com/labstack/echo/v4"

// A route definition we can use to wire up all routes
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc echo.HandlerFunc
}

type Routes []Route
