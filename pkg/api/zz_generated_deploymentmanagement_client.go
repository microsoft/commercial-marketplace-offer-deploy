//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package api

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DeploymentManagementClient contains the methods for the DeploymentManagementClient group.
// Don't use this type directly, use a constructor function instead.
type DeploymentManagementClient struct {
	internal *azcore.Client
	endpoint string
}

// CreatEventSubscription - Create a subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - body - Create event subscription
//   - options - DeploymentManagementClientCreatEventSubscriptionOptions contains the optional parameters for the DeploymentManagementClient.CreatEventSubscription
//     method.
func (client *DeploymentManagementClient) CreatEventSubscription(ctx context.Context, body CreateEventSubscription, options *DeploymentManagementClientCreatEventSubscriptionOptions) (DeploymentManagementClientCreatEventSubscriptionResponse, error) {
	req, err := client.creatEventSubscriptionCreateRequest(ctx, body, options)
	if err != nil {
		return DeploymentManagementClientCreatEventSubscriptionResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientCreatEventSubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusMethodNotAllowed) {
		return DeploymentManagementClientCreatEventSubscriptionResponse{}, runtime.NewResponseError(resp)
	}
	return client.creatEventSubscriptionHandleResponse(resp)
}

// creatEventSubscriptionCreateRequest creates the CreatEventSubscription request.
func (client *DeploymentManagementClient) creatEventSubscriptionCreateRequest(ctx context.Context, body CreateEventSubscription, options *DeploymentManagementClientCreatEventSubscriptionOptions) (*policy.Request, error) {
	urlPath := "/events/subscriptions"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// creatEventSubscriptionHandleResponse handles the CreatEventSubscription response.
func (client *DeploymentManagementClient) creatEventSubscriptionHandleResponse(resp *http.Response) (DeploymentManagementClientCreatEventSubscriptionResponse, error) {
	result := DeploymentManagementClientCreatEventSubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CreateEventSubscriptionResponse); err != nil {
		return DeploymentManagementClientCreatEventSubscriptionResponse{}, err
	}
	return result, nil
}

// CreateDeployment - Creates a new deployment instances
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - body - Deployment object that needs to be added to the store
//   - options - DeploymentManagementClientCreateDeploymentOptions contains the optional parameters for the DeploymentManagementClient.CreateDeployment
//     method.
func (client *DeploymentManagementClient) CreateDeployment(ctx context.Context, body CreateDeployment, options *DeploymentManagementClientCreateDeploymentOptions) (DeploymentManagementClientCreateDeploymentResponse, error) {
	req, err := client.createDeploymentCreateRequest(ctx, body, options)
	if err != nil {
		return DeploymentManagementClientCreateDeploymentResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientCreateDeploymentResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusMethodNotAllowed) {
		return DeploymentManagementClientCreateDeploymentResponse{}, runtime.NewResponseError(resp)
	}
	return client.createDeploymentHandleResponse(resp)
}

// createDeploymentCreateRequest creates the CreateDeployment request.
func (client *DeploymentManagementClient) createDeploymentCreateRequest(ctx context.Context, body CreateDeployment, options *DeploymentManagementClientCreateDeploymentOptions) (*policy.Request, error) {
	urlPath := "/deployments"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// createDeploymentHandleResponse handles the CreateDeployment response.
func (client *DeploymentManagementClient) createDeploymentHandleResponse(resp *http.Response) (DeploymentManagementClientCreateDeploymentResponse, error) {
	result := DeploymentManagementClientCreateDeploymentResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Deployment); err != nil {
		return DeploymentManagementClientCreateDeploymentResponse{}, err
	}
	return result, nil
}

// DeleteEventSubscription - Deletes a subscription to an even type
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - subscriptionID - ID of subscription
//   - options - DeploymentManagementClientDeleteEventSubscriptionOptions contains the optional parameters for the DeploymentManagementClient.DeleteEventSubscription
//     method.
func (client *DeploymentManagementClient) DeleteEventSubscription(ctx context.Context, subscriptionID string, options *DeploymentManagementClientDeleteEventSubscriptionOptions) (DeploymentManagementClientDeleteEventSubscriptionResponse, error) {
	req, err := client.deleteEventSubscriptionCreateRequest(ctx, subscriptionID, options)
	if err != nil {
		return DeploymentManagementClientDeleteEventSubscriptionResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientDeleteEventSubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusBadRequest, http.StatusNotFound) {
		return DeploymentManagementClientDeleteEventSubscriptionResponse{}, runtime.NewResponseError(resp)
	}
	return DeploymentManagementClientDeleteEventSubscriptionResponse{}, nil
}

// deleteEventSubscriptionCreateRequest creates the DeleteEventSubscription request.
func (client *DeploymentManagementClient) deleteEventSubscriptionCreateRequest(ctx context.Context, subscriptionID string, options *DeploymentManagementClientDeleteEventSubscriptionOptions) (*policy.Request, error) {
	urlPath := "/events/subscriptions/{subscriptionId}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	return req, nil
}

// GetDeployment - Returns a single deployment
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - deploymentID - ID of deployment to return
//   - options - DeploymentManagementClientGetDeploymentOptions contains the optional parameters for the DeploymentManagementClient.GetDeployment
//     method.
func (client *DeploymentManagementClient) GetDeployment(ctx context.Context, deploymentID int32, options *DeploymentManagementClientGetDeploymentOptions) (DeploymentManagementClientGetDeploymentResponse, error) {
	req, err := client.getDeploymentCreateRequest(ctx, deploymentID, options)
	if err != nil {
		return DeploymentManagementClientGetDeploymentResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientGetDeploymentResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusBadRequest, http.StatusNotFound) {
		return DeploymentManagementClientGetDeploymentResponse{}, runtime.NewResponseError(resp)
	}
	return client.getDeploymentHandleResponse(resp)
}

// getDeploymentCreateRequest creates the GetDeployment request.
func (client *DeploymentManagementClient) getDeploymentCreateRequest(ctx context.Context, deploymentID int32, options *DeploymentManagementClientGetDeploymentOptions) (*policy.Request, error) {
	urlPath := "/deployments/{deploymentId}"
	urlPath = strings.ReplaceAll(urlPath, "{deploymentId}", url.PathEscape(strconv.FormatInt(int64(deploymentID), 10)))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getDeploymentHandleResponse handles the GetDeployment response.
func (client *DeploymentManagementClient) getDeploymentHandleResponse(resp *http.Response) (DeploymentManagementClientGetDeploymentResponse, error) {
	result := DeploymentManagementClientGetDeploymentResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Deployment); err != nil {
		return DeploymentManagementClientGetDeploymentResponse{}, err
	}
	return result, nil
}

// GetEventSubscription - Gets a subscription to an even type
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - subscriptionID - ID of subscription
//   - options - DeploymentManagementClientGetEventSubscriptionOptions contains the optional parameters for the DeploymentManagementClient.GetEventSubscription
//     method.
func (client *DeploymentManagementClient) GetEventSubscription(ctx context.Context, subscriptionID string, options *DeploymentManagementClientGetEventSubscriptionOptions) (DeploymentManagementClientGetEventSubscriptionResponse, error) {
	req, err := client.getEventSubscriptionCreateRequest(ctx, subscriptionID, options)
	if err != nil {
		return DeploymentManagementClientGetEventSubscriptionResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientGetEventSubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusBadRequest, http.StatusNotFound) {
		return DeploymentManagementClientGetEventSubscriptionResponse{}, runtime.NewResponseError(resp)
	}
	return client.getEventSubscriptionHandleResponse(resp)
}

// getEventSubscriptionCreateRequest creates the GetEventSubscription request.
func (client *DeploymentManagementClient) getEventSubscriptionCreateRequest(ctx context.Context, subscriptionID string, options *DeploymentManagementClientGetEventSubscriptionOptions) (*policy.Request, error) {
	urlPath := "/events/subscriptions/{subscriptionId}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getEventSubscriptionHandleResponse handles the GetEventSubscription response.
func (client *DeploymentManagementClient) getEventSubscriptionHandleResponse(resp *http.Response) (DeploymentManagementClientGetEventSubscriptionResponse, error) {
	result := DeploymentManagementClientGetEventSubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventSubscription); err != nil {
		return DeploymentManagementClientGetEventSubscriptionResponse{}, err
	}
	return result, nil
}

// GetEventTypes - Returns a list of all event types
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - options - DeploymentManagementClientGetEventTypesOptions contains the optional parameters for the DeploymentManagementClient.GetEventTypes
//     method.
func (client *DeploymentManagementClient) GetEventTypes(ctx context.Context, options *DeploymentManagementClientGetEventTypesOptions) (DeploymentManagementClientGetEventTypesResponse, error) {
	req, err := client.getEventTypesCreateRequest(ctx, options)
	if err != nil {
		return DeploymentManagementClientGetEventTypesResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientGetEventTypesResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DeploymentManagementClientGetEventTypesResponse{}, runtime.NewResponseError(resp)
	}
	return client.getEventTypesHandleResponse(resp)
}

// getEventTypesCreateRequest creates the GetEventTypes request.
func (client *DeploymentManagementClient) getEventTypesCreateRequest(ctx context.Context, options *DeploymentManagementClientGetEventTypesOptions) (*policy.Request, error) {
	urlPath := "/events"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getEventTypesHandleResponse handles the GetEventTypes response.
func (client *DeploymentManagementClient) getEventTypesHandleResponse(resp *http.Response) (DeploymentManagementClientGetEventTypesResponse, error) {
	result := DeploymentManagementClientGetEventTypesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventTypeArray); err != nil {
		return DeploymentManagementClientGetEventTypesResponse{}, err
	}
	return result, nil
}

// GetInvokedDeploymentOperation - Gets the state of a command operation that's been invoked
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - operationID - ID of the triggered operation
//   - options - DeploymentManagementClientGetInvokedDeploymentOperationOptions contains the optional parameters for the DeploymentManagementClient.GetInvokedDeploymentOperation
//     method.
func (client *DeploymentManagementClient) GetInvokedDeploymentOperation(ctx context.Context, operationID string, options *DeploymentManagementClientGetInvokedDeploymentOperationOptions) (DeploymentManagementClientGetInvokedDeploymentOperationResponse, error) {
	req, err := client.getInvokedDeploymentOperationCreateRequest(ctx, operationID, options)
	if err != nil {
		return DeploymentManagementClientGetInvokedDeploymentOperationResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientGetInvokedDeploymentOperationResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DeploymentManagementClientGetInvokedDeploymentOperationResponse{}, runtime.NewResponseError(resp)
	}
	return client.getInvokedDeploymentOperationHandleResponse(resp)
}

// getInvokedDeploymentOperationCreateRequest creates the GetInvokedDeploymentOperation request.
func (client *DeploymentManagementClient) getInvokedDeploymentOperationCreateRequest(ctx context.Context, operationID string, options *DeploymentManagementClientGetInvokedDeploymentOperationOptions) (*policy.Request, error) {
	urlPath := "/deployments/operations/{operationId}"
	urlPath = strings.ReplaceAll(urlPath, "{operationId}", url.PathEscape(operationID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getInvokedDeploymentOperationHandleResponse handles the GetInvokedDeploymentOperation response.
func (client *DeploymentManagementClient) getInvokedDeploymentOperationHandleResponse(resp *http.Response) (DeploymentManagementClientGetInvokedDeploymentOperationResponse, error) {
	result := DeploymentManagementClientGetInvokedDeploymentOperationResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.InvokedDeploymentOperationResponse); err != nil {
		return DeploymentManagementClientGetInvokedDeploymentOperationResponse{}, err
	}
	return result, nil
}

// InvokeDeploymentOperation - Invokes a deployment operation with parameters
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - deploymentID - ID of deployment
//   - body - Deployment operation invocation
//   - options - DeploymentManagementClientInvokeDeploymentOperationOptions contains the optional parameters for the DeploymentManagementClient.InvokeDeploymentOperation
//     method.
func (client *DeploymentManagementClient) InvokeDeploymentOperation(ctx context.Context, deploymentID int32, body InvokeDeploymentOperationRequest, options *DeploymentManagementClientInvokeDeploymentOperationOptions) (DeploymentManagementClientInvokeDeploymentOperationResponse, error) {
	req, err := client.invokeDeploymentOperationCreateRequest(ctx, deploymentID, body, options)
	if err != nil {
		return DeploymentManagementClientInvokeDeploymentOperationResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientInvokeDeploymentOperationResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DeploymentManagementClientInvokeDeploymentOperationResponse{}, runtime.NewResponseError(resp)
	}
	return client.invokeDeploymentOperationHandleResponse(resp)
}

// invokeDeploymentOperationCreateRequest creates the InvokeDeploymentOperation request.
func (client *DeploymentManagementClient) invokeDeploymentOperationCreateRequest(ctx context.Context, deploymentID int32, body InvokeDeploymentOperationRequest, options *DeploymentManagementClientInvokeDeploymentOperationOptions) (*policy.Request, error) {
	urlPath := "/deployment/{deploymentId}/operation"
	urlPath = strings.ReplaceAll(urlPath, "{deploymentId}", url.PathEscape(strconv.FormatInt(int64(deploymentID), 10)))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// invokeDeploymentOperationHandleResponse handles the InvokeDeploymentOperation response.
func (client *DeploymentManagementClient) invokeDeploymentOperationHandleResponse(resp *http.Response) (DeploymentManagementClientInvokeDeploymentOperationResponse, error) {
	result := DeploymentManagementClientInvokeDeploymentOperationResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.InvokedDeploymentOperationResponse); err != nil {
		return DeploymentManagementClientInvokeDeploymentOperationResponse{}, err
	}
	return result, nil
}

// ListDeployments - List all deployments
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - options - DeploymentManagementClientListDeploymentsOptions contains the optional parameters for the DeploymentManagementClient.ListDeployments
//     method.
func (client *DeploymentManagementClient) ListDeployments(ctx context.Context, options *DeploymentManagementClientListDeploymentsOptions) (DeploymentManagementClientListDeploymentsResponse, error) {
	req, err := client.listDeploymentsCreateRequest(ctx, options)
	if err != nil {
		return DeploymentManagementClientListDeploymentsResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientListDeploymentsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DeploymentManagementClientListDeploymentsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listDeploymentsHandleResponse(resp)
}

// listDeploymentsCreateRequest creates the ListDeployments request.
func (client *DeploymentManagementClient) listDeploymentsCreateRequest(ctx context.Context, options *DeploymentManagementClientListDeploymentsOptions) (*policy.Request, error) {
	urlPath := "/deployments"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Status != nil {
			for _, qv := range options.Status {
		reqQP.Add("status", string(qv))
	}
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listDeploymentsHandleResponse handles the ListDeployments response.
func (client *DeploymentManagementClient) listDeploymentsHandleResponse(resp *http.Response) (DeploymentManagementClientListDeploymentsResponse, error) {
	result := DeploymentManagementClientListDeploymentsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeploymentArray); err != nil {
		return DeploymentManagementClientListDeploymentsResponse{}, err
	}
	return result, nil
}

// ListEventSubscriptions - List all subscriptions
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - options - DeploymentManagementClientListEventSubscriptionsOptions contains the optional parameters for the DeploymentManagementClient.ListEventSubscriptions
//     method.
func (client *DeploymentManagementClient) ListEventSubscriptions(ctx context.Context, options *DeploymentManagementClientListEventSubscriptionsOptions) (DeploymentManagementClientListEventSubscriptionsResponse, error) {
	req, err := client.listEventSubscriptionsCreateRequest(ctx, options)
	if err != nil {
		return DeploymentManagementClientListEventSubscriptionsResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientListEventSubscriptionsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusBadRequest, http.StatusNotFound) {
		return DeploymentManagementClientListEventSubscriptionsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listEventSubscriptionsHandleResponse(resp)
}

// listEventSubscriptionsCreateRequest creates the ListEventSubscriptions request.
func (client *DeploymentManagementClient) listEventSubscriptionsCreateRequest(ctx context.Context, options *DeploymentManagementClientListEventSubscriptionsOptions) (*policy.Request, error) {
	urlPath := "/events/subscriptions"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listEventSubscriptionsHandleResponse handles the ListEventSubscriptions response.
func (client *DeploymentManagementClient) listEventSubscriptionsHandleResponse(resp *http.Response) (DeploymentManagementClientListEventSubscriptionsResponse, error) {
	result := DeploymentManagementClientListEventSubscriptionsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventSubscriptionArray); err != nil {
		return DeploymentManagementClientListEventSubscriptionsResponse{}, err
	}
	return result, nil
}

// ListOperations - Returns a list of available operations that can be performed on a deployment
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - options - DeploymentManagementClientListOperationsOptions contains the optional parameters for the DeploymentManagementClient.ListOperations
//     method.
func (client *DeploymentManagementClient) ListOperations(ctx context.Context, options *DeploymentManagementClientListOperationsOptions) (DeploymentManagementClientListOperationsResponse, error) {
	req, err := client.listOperationsCreateRequest(ctx, options)
	if err != nil {
		return DeploymentManagementClientListOperationsResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientListOperationsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DeploymentManagementClientListOperationsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listOperationsHandleResponse(resp)
}

// listOperationsCreateRequest creates the ListOperations request.
func (client *DeploymentManagementClient) listOperationsCreateRequest(ctx context.Context, options *DeploymentManagementClientListOperationsOptions) (*policy.Request, error) {
	urlPath := "/operations"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listOperationsHandleResponse handles the ListOperations response.
func (client *DeploymentManagementClient) listOperationsHandleResponse(resp *http.Response) (DeploymentManagementClientListOperationsResponse, error) {
	result := DeploymentManagementClientListOperationsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.OperationArray); err != nil {
		return DeploymentManagementClientListOperationsResponse{}, err
	}
	return result, nil
}

// UpdateDeployment - Update an existing deployment
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 0.1.0
//   - body - Deployment object that needs to be added to the store
//   - options - DeploymentManagementClientUpdateDeploymentOptions contains the optional parameters for the DeploymentManagementClient.UpdateDeployment
//     method.
func (client *DeploymentManagementClient) UpdateDeployment(ctx context.Context, body Deployment, options *DeploymentManagementClientUpdateDeploymentOptions) (DeploymentManagementClientUpdateDeploymentResponse, error) {
	req, err := client.updateDeploymentCreateRequest(ctx, body, options)
	if err != nil {
		return DeploymentManagementClientUpdateDeploymentResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeploymentManagementClientUpdateDeploymentResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusBadRequest, http.StatusNotFound, http.StatusMethodNotAllowed) {
		return DeploymentManagementClientUpdateDeploymentResponse{}, runtime.NewResponseError(resp)
	}
	return DeploymentManagementClientUpdateDeploymentResponse{}, nil
}

// updateDeploymentCreateRequest creates the UpdateDeployment request.
func (client *DeploymentManagementClient) updateDeploymentCreateRequest(ctx context.Context, body Deployment, options *DeploymentManagementClientUpdateDeploymentOptions) (*policy.Request, error) {
	urlPath := "/deployments"
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	return req, runtime.MarshalAsJSON(req, body)
}

