package sdk

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
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
