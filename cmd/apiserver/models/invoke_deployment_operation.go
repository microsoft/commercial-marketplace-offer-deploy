package models

type InvokeDeploymentOperation struct {
	Name string `json:"name,omitempty"`

	Parameters []OperationParameter `json:"parameters,omitempty"`
}
