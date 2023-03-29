package eventgrid

import (
	"fmt"
	"net/http"
	"strings"

	"gobn.github.io/coalesce"
)

const (
	// the header value for the event grid validation request
	HeaderNameEventType          = "aeg-event-type"
	RequestSubscriptionEventType = "SubscriptionValidation"
	ValidationEventType          = "Microsoft.EventGrid.SubscriptionValidationEvent"
)

// Bind binds the request body into provided type `to`.
type RequestBodyBinderFunc func(to any) error

// the web hook endpoint validation response body that should be returned by an HTTP handler performing the authentication handshake
type WebHookValidationEventResponse struct {
	ValidationResponse string `json:"validationResponse"`
}

// the validation event result
type WebHookValidationEventResult struct {
	// whether the request was handled. If false, then the http request wasn't a validation event handshake request
	Handled  bool
	Error    error
	Response *WebHookValidationEventResponse
}

// handler for web hook validation event
type WebHookValidationEventHandler interface {
	// performs a check if the http request is a a validation event. If it is, it attempts to
	// create a validation response and return it
	//	returns: the validation response and an error indicating whether the request was identified as a validation event
	Handle(context *http.Request) *WebHookValidationEventResult
}

// validates webhook endpoint for sync handshakes
//
//	see: https://learn.microsoft.com/en-us/azure/event-grid/webhook-event-delivery for more details
type webHookEndpointValidator struct {
	BinderFunc RequestBodyBinderFunc
	Result     *WebHookValidationEventResult
}

// Handle implements WebHookValidationEventHandler
func (v *webHookEndpointValidator) Handle(request *http.Request) *WebHookValidationEventResult {
	result := v.Result
	result.Handled = false

	if v.isValidationRequest(&request.Header) {
		result.Handled = true
		if request, err := v.getValidationRequest(); err == nil {
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

	return result
}

func NewWebHookValidationEventHandler(binderFunc RequestBodyBinderFunc) WebHookValidationEventHandler {
	return &webHookEndpointValidator{
		BinderFunc: binderFunc,
		Result:     &WebHookValidationEventResult{},
	}
}

type validationRequest struct {
	Data struct {
		ValidationCode string `json:"validationCode"`
	} `json:"data"`
	EventType string `json:"eventType"`
}

func (v *webHookEndpointValidator) getValidationRequest() (*validationRequest, error) {
	// the payload will containe just 1 validation element in an array
	payload := []validationRequest{}
	err := v.BinderFunc(&payload)

	if err != nil || len(payload) != 1 {
		return nil, coalesce.Error(err, fmt.Errorf("validation failed"))
	}
	return &payload[0], nil
}

// identifies if this is a validation request that's used to handshake with EventGrid
//
//	see: https://learn.microsoft.com/en-us/azure/event-grid/webhook-event-delivery for more information
func (v *webHookEndpointValidator) isValidationRequest(header *http.Header) bool {
	value := header.Get(HeaderNameEventType)
	return value != "" && strings.EqualFold(value, RequestSubscriptionEventType)
}
