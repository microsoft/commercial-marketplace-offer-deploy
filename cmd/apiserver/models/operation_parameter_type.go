package models

// The parameter type information for a parameter of an operation
type OperationParameterType struct {
	Name string `json:"name,omitempty"`

	Type_ string `json:"type,omitempty"`
}
