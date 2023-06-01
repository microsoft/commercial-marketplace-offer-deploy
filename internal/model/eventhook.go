package model

type EventHook struct {
	BaseWithGuidPrimaryKey
	Callback string `json:"callback"`
	Name     string `json:"name" gorm:"unique"`
	ApiKey   string `json:"authKey"`
}
