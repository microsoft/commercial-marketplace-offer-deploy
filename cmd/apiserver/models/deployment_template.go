package models

type DeploymentTemplate struct {
	Name string `json:"name,omitempty"`

	Uri string `json:"uri,omitempty"`

	Parameters []DeploymentTemplateParameter `json:"parameters,omitempty"`
}
