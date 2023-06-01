package model

type Stage struct {
	BaseWithGuidPrimaryKey
	DeploymentID        uint   `json:"deploymentId"`
	Name                string `json:"name"`
	AzureDeploymentName string `json:"azureDeploymentName"`

	// the default number of retries
	Retries uint `json:"retries"`
}
