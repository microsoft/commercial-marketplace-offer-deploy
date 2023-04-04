package config

import (
	"os"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/app"
	"github.com/stretchr/testify/assert"
)

func TestGetAppSettings(t *testing.T) {

	// add some env variables for testing to make sure they are loaded
	os.Setenv("AZURE_CLIENT_ID", "test-client-id")
	os.Setenv("AZURE_TENANT_ID", "test-tenant-id")
	os.Setenv("AZURE_SUBSCRIPTION_ID", "test-subscription-id")
	os.Setenv("AZURE_RESOURCE_GROUP", "test-resource-group")
	os.Setenv("AZURE_LOCATION", "test-location")
	os.Setenv("AZURE_SERVICEBUS_NAMESPACE", "test-servicebus-namespace")

	app.GetApp("")

	result := GetAppSettings()

	if assert.NotNil(t, result) {
		assert.Equal(t, "test-client-id", result.Azure.ClientId)
		assert.Equal(t, "test-tenant-id", result.Azure.TenantId)
		assert.Equal(t, "test-subscription-id", result.Azure.SubscriptionId)
		assert.Equal(t, "test-resource-group", result.Azure.ResourceGroupName)
		assert.Equal(t, "test-location", result.Azure.Location)
		assert.Equal(t, "test-servicebus-namespace", result.Azure.ServiceBusNamespace)
	}

}
