package subscriptionvalidation

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gobn.github.io/coalesce"
)

const (
	// the header value for the event grid validation request
	HeaderNameEventType          = "aeg-event-type"
	RequestSubscriptionEventType = "SubscriptionValidation"
	ValidationEventType          = "Microsoft.EventGrid.SubscriptionValidationEvent"
)

// validates and EventGrid webhook subscription for synchronous handshakes with Azure Event Grid
//
//	see: https://learn.microsoft.com/en-us/azure/event-grid/webhook-event-delivery for more details
type WebHookValidationEventClient interface {
	// performs a check if the http request is a a validation event. If it is, it attempts to
	// create a validation response and return it
	//	returns: the validation response and an error indicating whether the request was identified as a validation event
	Validate() *WebHookValidationEventResult
}

// internal implementation of WebHookValidationEventClient
type client struct {
	context echo.Context
}

// create a new web hook endpoint client
func NewWebHookValidationEventClient(c echo.Context) WebHookValidationEventClient {
	return &client{context: c}
}

// Handle implements WebHookValidationEventHandler
func (c *client) Validate() *WebHookValidationEventResult {
	result := WebHookValidationEventResult{Handled: false}
	request := c.context.Request()

	if c.isValidationRequest(&request.Header) {
		result.Handled = true
		if request, err := c.getValidationRequest(); err == nil {
			if request.EventType == ValidationEventType {
				result.Response = &WebHookValidationEventResponse{
					ValidationResponse: request.Data.ValidationCode,
				}
			} else {
				result.Error = fmt.Errorf("")
			}
		} else {
			result.Error = err
		}
	}

	return &result
}

// gets the validation request from the http request
func (c *client) getValidationRequest() (*validationRequest, error) {
	// the payload will containe just 1 validation element in an array
	payload := []validationRequest{}
	err := c.context.Bind(&payload)

	if err != nil || len(payload) != 1 {
		return nil, coalesce.Error(err, fmt.Errorf("validation failed"))
	}
	return &payload[0], nil
}

// identifies if this is a validation request that's used to handshake with EventGrid
//
//	see: https://learn.microsoft.com/en-us/azure/event-grid/webhook-event-delivery for more information
func (c *client) isValidationRequest(header *http.Header) bool {
	value := header.Get(HeaderNameEventType)
	return value != "" && strings.EqualFold(value, RequestSubscriptionEventType)
}
