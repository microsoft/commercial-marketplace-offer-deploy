package models

// Defines an available operation
type Operation struct {
	Target *OperationTargetType `json:"target,omitempty"`

	Name string `json:"name,omitempty"`

	Parameters []OperationParameterType `json:"parameters,omitempty"`
}
