package models

type OperationParameter struct {
	Name string `json:"name,omitempty"`

	Value string `json:"value,omitempty"`

	Type_ string `json:"type,omitempty"`
}
