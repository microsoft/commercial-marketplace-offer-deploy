package sdk

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("test", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	cloudConfig := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "test", Services: map[cloud.ServiceName]cloud.ServiceConfiguration{},
	}
	_, err = NewClient("test", nil, &ClientOptions{ClientOptions: azcore.ClientOptions{Cloud: cloudConfig}})
	require.NoErrorf(t, err, "No error.")
}

func TestListDeployments(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	client, err := NewClient("http://localhost:8080", cred, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	result, err := client.ListDeployments(context.TODO())

	if err != nil {
		t.Logf("Error: %s", err)
	}
	require.NotNil(t, result)
}
