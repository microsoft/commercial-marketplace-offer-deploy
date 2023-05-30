package model

type Stage struct {
	BaseWithGuidPrimaryKey
	Name           string `json:"name"`
	DeploymentName string `json:"deploymentName"`

	// the default number of retries
	Retries uint `json:"retries"`
}

type EventHook struct {
	BaseWithGuidPrimaryKey
	Callback string `json:"callback"`
	Name     string `json:"name" gorm:"unique"`
	ApiKey   string `json:"authKey"`
}
