package app

import (
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/stretchr/testify/assert"
)

func Test_createOptions_url(t *testing.T) {
	appConfig := &config.AppConfig{
		Http: config.HttpSettings{
			BaseUrl: "https://localhost:8080",
		},
	}
	result := createOptions(appConfig)

	assert.Equal(t, "https://localhost:8080/eventgrid", result.EndpointUrl)
}

func Test_getSubscriptionName(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "Microsoft-test_rg", want: "testrg-events-" + getHostname()},
		{input: "EventGrid-test_rg", want: "testrg-events-" + getHostname()},
		{input: "System-test_rg", want: "testrg-events-" + getHostname()},
		{input: "rg-that-exceeds-the-max-length-of-sixty-four-characters-long", want: "rg-that-exceeds-the-max-events-" + getHostname()},
		{input: ",.~`{}|/<>[]rg-with-special-*&^%$#@!_+=.:'\"", want: "rg-with-special-events-" + getHostname()},
	}

	for _, tc := range tests {
		got := getSubscriptionName(tc.input)
		assert.Equal(t, tc.want, got)
	}
}
