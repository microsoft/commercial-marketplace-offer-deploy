package sdk

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/require"
)

const testEndPoint = "http://localhost:8080"

func TestStartDeployment(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	client, err := NewClient(testEndPoint, cred, nil)

	require.NoError(t, err)
	require.NotNil(t, client)

	// TODO: construct the deployment
	result, err := client.StartDeployment(context.TODO(), 123, nil)

	if err != nil {
		t.Logf("Error: %s", err)
	}
	require.NotNil(t, result)
}
