package api

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func NewDeploymentManagementClient(endpoint string, credential azcore.TokenCredential, options *policy.ClientOptions) (*DeploymentManagementClient, error) {
	tokenScope, err := getDefaultScope(endpoint)

	if err != nil {
		return nil, err
	}

	internalClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(credential, []string{tokenScope}, nil),
		},
	}, options)

	if err != nil {
		return nil, err
	}

	deploymentClient := &DeploymentManagementClient{
		internal: internalClient,
		endpoint: endpoint,
	}

	return deploymentClient, nil
}

func getDefaultScope(endpoint string) (string, error) {
	return "api://modm/.default", nil
}