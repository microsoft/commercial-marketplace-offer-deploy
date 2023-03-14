package models

type InvokedOperationTarget struct {
	Type_ string `json:"type,omitempty"`

	Id *InvokedOperationTargetId `json:"id,omitempty"`
}
