package api

import (
	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type DeploymentManagementClientOptions struct {
	ClientName    string
	Version       string
	ClientOptions *policy.ClientOptions
}

func NewDeploymentManagementClient(endpoint string, credential azcore.TokenCredential, options *DeploymentManagementClientOptions) (*DeploymentManagementClient, error) {
	tokenScope, err := getDefaultScope(endpoint)

	if err != nil {
		return nil, err
	}

	log.Debugf("options is %v", options)

	internalClient, err := azcore.NewClient(options.ClientName, options.Version, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(credential, []string{tokenScope}, nil),
		},
	}, options.ClientOptions)

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
	return "https://management.azure.com/.default", nil
}
