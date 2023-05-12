package sdk

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
)

// DefaultRetries is the default number of retries for an operation against a deployment
// default is 3
const DefaultRetries = 3

// Performs a dry run of a deployment and returns the verification results
// returns: verification results
func (client *Client) DryRun(ctx context.Context, deploymentId int, templateParameters map[string]interface{}) (*DryRunResponse, error) {
	retries := DefaultRetries
	response, err := client.invokeDeploymentOperation(ctx, true, operation.TypeDryRun, deploymentId, templateParameters, retries)
	if err != nil {
		return nil, err
	}
	return &DryRunResponse{
		Id:      uuid.MustParse(*response.ID),
		Results: response.Result,
		Status:  *response.Status,
	}, nil
}

// Starts a deployment, asynchonously. The deployment response immediately returns with a status of "accepted" or "rejected".
//
//	remarks: The deployment will be executed asynchronously, and the status of the deployment can be queried using the GetDeployment operation.
//	id: deployment id
//
// returns: the unique UUID of the deployment execution instance
func (client *Client) Start(ctx context.Context, deploymentId int, templateParameters map[string]interface{}, options *StartOptions) (*StartDeploymentResponse, error) {
	retries := DefaultRetries
	if options != nil {
		retries = options.Retries
	}
	response, err := client.invokeDeploymentOperation(ctx, false, operation.TypeStartDeployment, deploymentId, templateParameters, retries)
	if err != nil {
		return nil, err
	}
	return &StartDeploymentResponse{
		Id:     uuid.MustParse(*response.ID),
		Status: *response.Status,
	}, nil
}

// Retries a deployment, regardless of the current status
//
//	id: deployment id
func (client *Client) Retry(ctx context.Context, deploymentId int, options *RetryOptions) (*RetryResponse, error) {
	operationType := operation.TypeRetryDeployment

	// if we have a stageId set, then we want to retry a stage of the deployment
	if options != nil {
		if options.StageId != uuid.Nil {
			operationType = operation.TypeRetryStage
		}
	}

	params := make(map[string]interface{})
	params["stageId"] = options.StageId

	retries := 0 // we don't want to retry a retry.
	resp, err := client.invokeDeploymentOperation(ctx, false, operationType, deploymentId, params, retries)
	if err != nil {
		return nil, err
	}

	return &RetryResponse{
		Id:         uuid.MustParse(*resp.ID),
		Status:     *resp.Status,
		Parameters: resp.Parameters.(map[string]any),
	}, nil
}

// Creates a deployment record that will be used to kick off all available deployment operations (dry run, start, retry, etc.)
func (client *Client) Create(ctx context.Context, request api.CreateDeployment) (*api.Deployment, error) {
	response, err := client.internalClient.CreateDeployment(ctx, request, nil)

	if err != nil {
		return nil, err
	}
	deployment := response.Deployment
	return &deployment, nil
}

func (client *Client) Get(ctx context.Context, deploymentId int) (*GetResponse, error) {
	resp, err := client.internalClient.GetDeployment(ctx, int32(deploymentId), nil)
	if err != nil {
		return nil, err
	}
	return &GetResponse{
		Deployment: &resp.Deployment,
	}, nil
}

// list all deployments
func (client *Client) List(ctx context.Context) (*ListResponse, error) {
	resp, err := client.internalClient.ListDeployments(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &ListResponse{
		Deployments: resp.DeploymentArray,
	}, nil
}
