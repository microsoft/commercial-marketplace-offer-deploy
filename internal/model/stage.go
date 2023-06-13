package model

type Stage struct {
	BaseWithGuidPrimaryKey
	DeploymentID        uint              `json:"deploymentId" gorm:"primaryKey"`
	Name                string            `json:"name"`
	AzureDeploymentName string            `json:"azureDeploymentName"`
	Attributes          map[string]string `json:"attributes" gorm:"json"`
	// the default number of retries
	Retries uint `json:"retries"`
}
