package subscriptionvalidation

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

// Private types

type validationRequest struct {
	Data struct {
		ValidationCode string `json:"validationCode"`
	} `json:"data"`
	EventType string `json:"eventType"`
}
