package eventgrid

// ResourceEventData is the data structure for the event grid event
// use only for unmarshalling in order to map to resource
type ResourceEventData struct {
	CorrelationID    string `json:"correlationId"`
	ResourceProvider string `json:"resourceProvider"`
	ResourceURI      string `json:"resourceUri"`
	OperationName    string `json:"operationName"`
	Status           string `json:"status"`
	SubscriptionID   string `json:"subscriptionId"`
	TenantID         string `json:"tenantId"`
}
