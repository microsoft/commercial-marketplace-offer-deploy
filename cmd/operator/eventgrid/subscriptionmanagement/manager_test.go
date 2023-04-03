package subscriptionmanagement

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: convert these test to integration tests that use https://github.com/ngrok/ngrok-go to tunnel to localhost
func TestCreateEventSubscription(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	resourceGroupId := "/subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/sandbox1"
	client, err := NewEventGridManager(cred, resourceGroupId)

	require.NoError(t, err)

	ctx := context.TODO()
	subscriptionName := "test-subscription"
	endpointUrl := "https://2e75-172-73-181-161.ngrok.io/eventgrid"

	result, err := client.CreateEventSubscription(ctx, subscriptionName, endpointUrl)

	assert.NotNil(t, result)
	require.NoError(t, err)

	assert.Equal(t, subscriptionName, *result.Name)
}

// integration test. Needs to run against a subscription
func TestDeploymentEventsClient(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	resourceGroupId := "/subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/sandbox1"
	client, err := NewEventGridManager(cred, resourceGroupId)

	require.NoError(t, err)

	result, err := client.CreateSystemTopic(context.TODO())
	log.Printf("result:\n %v", result)

	require.NoError(t, err)
	assert.NotNil(t, result)

}
