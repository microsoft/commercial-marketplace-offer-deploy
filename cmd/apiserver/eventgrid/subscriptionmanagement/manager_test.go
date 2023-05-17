package subscriptionmanagement

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetSystemTopicName(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "Microsoft-test_rg", want: "testrg-event-topic"},
		{input: "EventGrid-test_rg", want: "testrg-event-topic"},
		{input: "System-test_rg", want: "testrg-event-topic"},
		{input: "rg-that-exceeds-the-max-length-of-fifty-characters-long", want: "rg-that-exceeds-the-max-length-of-fift-event-topic"},
		{input: ",.~`{}|/<>[]rg-with-special-*&^%$#@!_+=.:'\"", want: "rg-with-special-event-topic"},
	}

	for _, tc := range tests {
		got := getSystemTopicName(tc.input)
		assert.Equal(t, tc.want, got)
	}
}
