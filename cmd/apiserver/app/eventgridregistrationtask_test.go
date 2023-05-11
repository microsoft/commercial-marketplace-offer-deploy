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
