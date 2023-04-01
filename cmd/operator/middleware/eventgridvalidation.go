package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid/subscriptionvalidation"
)

// EventGridSubscriptionValidationMiddleware is a middleware that validates the event grid subscription
// It identifies the event grid subscription validation request. If the request is a validation request
// it returns the validation response. If the request is not a validation request, it continues to the next handler.
func EventGridWebHookSubscriptionValidation() echo.MiddlewareFunc {
	return eventGridSubscriptionValidationHandler
}

// the middleware handler
func eventGridSubscriptionValidationHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Print("Validating event grid subscription")
		validationResult := getResult(c)
		if validationResult != nil {
			return validationResult
		}
		return next(c)
	}
}

func getResult(c echo.Context) error {
	webhookValidator := subscriptionvalidation.NewWebHookValidationEventClient(c)
	result := webhookValidator.Validate()

	if result.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
	}

	if result.Handled {
		return c.JSON(http.StatusOK, &result.Response)
	}

	return nil
}
