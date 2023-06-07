package model

import "gorm.io/gorm"

// event hook subscriptions are stored with this
type EventHook struct {
	BaseWithGuidPrimaryKey
	Callback string `json:"callback"`
	Name     string `json:"name" gorm:"unique"`
	ApiKey   string `json:"authKey"`
}

// record for all event hook messages sent
type EventHookAuditEntry struct {
	gorm.Model
	Hash  	string `json:"hash"`
	Message any    `json:"message" gorm:"json"`
}
