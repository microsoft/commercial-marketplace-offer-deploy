package sdk

import (
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
)

// Client is the struct for interacting with an Azure App Configuration instance.
type Client struct {
	internalClient *generated.DeploymentManagementClient
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

	tokenScope, err := getDefaultScope(endpoint)

	if err != nil {
		return nil, err
	}

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(credential, []string{tokenScope}, nil),
		},
	}, &options.ClientOptions)

	internalClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(credential, []string{tokenScope}, nil),
		},
	}, &options.ClientOptions)

	deploymentClient := &generated.DeploymentManagementClient{
		internal: internalClient,
		endpoint: endpoint,
	}

	return &Client{internalClient: deploymentClient}, nil
}

func getDefaultScope(endpoint string) (string, error) {
	return "api://modm/.default", nil
}
