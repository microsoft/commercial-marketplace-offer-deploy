package models

type CreateEventSubscription struct {
	Name string `json:"name,omitempty"`

	Callback string `json:"callback,omitempty"`
}
