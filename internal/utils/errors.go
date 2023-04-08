package utils

import (
	"errors"
	"fmt"
)

// Joins the error messages into a single error
// Reason:
//
//	1.18 doesn't have errors.Join. This facilicates easy aggregate errors
func NewAggregateError(messages []string) error {
	if len(messages) == 0 {
		return nil
	}
	aggregate := errors.New("aggregate error")
	for index, message := range messages {
		aggregate = fmt.Errorf("error %d: %w", index, errors.New(message))
	}
	return aggregate
}
