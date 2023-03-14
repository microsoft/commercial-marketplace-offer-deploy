package models

type Deployment struct {
	Id int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Status string `json:"status,omitempty"`

	Template *DeploymentTemplate `json:"template,omitempty"`
}
