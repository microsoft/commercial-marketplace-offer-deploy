//go:build integration

package eventgrid

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// integration test. Needs to run against a subscription
// example call: go test -timeout 30s -run ^TestDeploymentEventsClient$ -tags=integration
func TestCreateEventSubscription(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	resourceGroupId := "/subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/sandbox1"
	client, err := NewDeploymentEventsClient(cred, resourceGroupId)

	require.NoError(t, err)

	result, err := client.CreateEventSubscription(context.TODO())
	log.Printf("result:\n %v", result)

	require.NoError(t, err)
	assert.NotNil(t, result)

}

// integration test. Needs to run against a subscription
func TestDeploymentEventsClient(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	resourceGroupId := "/subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/sandbox1"
	client, err := NewDeploymentEventsClient(cred, resourceGroupId)

	require.NoError(t, err)

	result, err := client.CreateSystemTopic(context.TODO())
	log.Printf("result:\n %v", result)

	require.NoError(t, err)
	assert.NotNil(t, result)

}
