package sdk

import (
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

// Client is the struct for interacting with an Azure App Configuration instance.
type Client struct {
	internalClient *DeploymentManagementClient
}

// ClientOptions contains the optional parameters for the NewClient method.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a new instance of Client with the specified values.
//   - endpoint - the endpoint of the Marketplace Offer Deployment Management service.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}

	internalClient, err := NewDeploymentManagementClient(endpoint, credential, &DeploymentManagementClientOptions{
		ClientOptions: &options.ClientOptions,
		ClientName:    moduleName + ".client",
		Version:       moduleVersion,
	})
	if err != nil {
		return nil, err
	}

	return &Client{internalClient: internalClient}, nil
}
