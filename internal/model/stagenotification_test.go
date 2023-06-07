package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StageNotification_AllSent(t *testing.T) {
	notification := &StageNotification{
		IsDone: false,
		Entries: []StageNotificationEntry{
			{IsSent: true},
		},
	}
	assert.True(t, notification.AllSent())
}
