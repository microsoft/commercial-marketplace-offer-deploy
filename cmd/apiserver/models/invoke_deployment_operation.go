package models

type InvokeDeploymentOperation struct {
	Name string `json:"name,omitempty"`
	Parameters map[string] interface{} `json:"parameters,omitempty"`
	//Parameters []OperationParameter `json:"parameters,omitempty"`
}
