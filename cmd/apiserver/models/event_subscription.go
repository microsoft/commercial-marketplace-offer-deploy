package models

type EventSubscription struct {
	Id string `json:"id,omitempty"`

	Topic string `json:"topic,omitempty"`

	Name string `json:"name,omitempty"`

	Callback string `json:"callback,omitempty"`
}
