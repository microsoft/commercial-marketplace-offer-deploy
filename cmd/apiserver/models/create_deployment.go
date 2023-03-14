package models

type CreateDeployment struct {
	Name string `json:"name"`

	MultiStage bool `json:"multiStage,omitempty"`

	Template *DeploymentTemplate `json:"template"`
}
