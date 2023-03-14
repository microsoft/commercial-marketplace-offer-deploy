package models

import (
	"time"
)

type InvokedOperation struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Target *InvokedOperationTarget `json:"target,omitempty"`

	Parameters []OperationParameter `json:"parameters,omitempty"`

	InvokedOn time.Time `json:"invokedOn,omitempty"`

	Status string `json:"status,omitempty"`
}
