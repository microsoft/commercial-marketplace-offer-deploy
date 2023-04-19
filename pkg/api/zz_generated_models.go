//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package api

import "time"

type CreateDeployment struct {
	// REQUIRED
	Name *string `json:"name,omitempty"`

	// REQUIRED; Anything
	Template any `json:"template,omitempty"`
	Location *string `json:"location,omitempty"`
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	SubscriptionID *string `json:"subscriptionId,omitempty"`
}

type CreateEventSubscriptionRequest struct {
	// API key to be used in the Authorization header, e.g. 'ApiKey =234dfsdf324234', to call the webhook callback URL.
	APIKey *string `json:"ApiKey,omitempty"`

	// The webhook callback
	Callback *string `json:"callback,omitempty"`

	// the name of the subscription
	Name *string `json:"name,omitempty"`
}

type CreateEventSubscriptionResponse struct {
	ID *string `json:"id,omitempty"`

	// the name of the subscription
	Name *string `json:"name,omitempty"`
}

type Deployment struct {
	ID *int32 `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Stages []*DeploymentStage `json:"stages,omitempty"`
	Status *string `json:"status,omitempty"`

	// Anything
	Template any `json:"template,omitempty"`
}

// DeploymentManagementClientCreatEventSubscriptionOptions contains the optional parameters for the DeploymentManagementClient.CreatEventSubscription
// method.
type DeploymentManagementClientCreatEventSubscriptionOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientCreateDeploymentOptions contains the optional parameters for the DeploymentManagementClient.CreateDeployment
// method.
type DeploymentManagementClientCreateDeploymentOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientDeleteEventSubscriptionOptions contains the optional parameters for the DeploymentManagementClient.DeleteEventSubscription
// method.
type DeploymentManagementClientDeleteEventSubscriptionOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientGetDeploymentOptions contains the optional parameters for the DeploymentManagementClient.GetDeployment
// method.
type DeploymentManagementClientGetDeploymentOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientGetEventSubscriptionOptions contains the optional parameters for the DeploymentManagementClient.GetEventSubscription
// method.
type DeploymentManagementClientGetEventSubscriptionOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientGetEventTypesOptions contains the optional parameters for the DeploymentManagementClient.GetEventTypes
// method.
type DeploymentManagementClientGetEventTypesOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientGetInvokedDeploymentOperationOptions contains the optional parameters for the DeploymentManagementClient.GetInvokedDeploymentOperation
// method.
type DeploymentManagementClientGetInvokedDeploymentOperationOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientInvokeDeploymentOperationOptions contains the optional parameters for the DeploymentManagementClient.InvokeDeploymentOperation
// method.
type DeploymentManagementClientInvokeDeploymentOperationOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientListDeploymentsOptions contains the optional parameters for the DeploymentManagementClient.ListDeployments
// method.
type DeploymentManagementClientListDeploymentsOptions struct {
	// Status values that need to be considered for filter
	Status []Status
}

// DeploymentManagementClientListEventSubscriptionsOptions contains the optional parameters for the DeploymentManagementClient.ListEventSubscriptions
// method.
type DeploymentManagementClientListEventSubscriptionsOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientListOperationsOptions contains the optional parameters for the DeploymentManagementClient.ListOperations
// method.
type DeploymentManagementClientListOperationsOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientUpdateDeploymentOptions contains the optional parameters for the DeploymentManagementClient.UpdateDeployment
// method.
type DeploymentManagementClientUpdateDeploymentOptions struct {
	// placeholder for future optional parameters
}

type DeploymentStage struct {
	ID *int32 `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
}

type EventSubscriptionResponse struct {
	// The webhook callback
	Callback *string `json:"callback,omitempty"`
	ID *string `json:"id,omitempty"`

	// the name of the subscription
	Name *string `json:"name,omitempty"`
}

type EventType struct {
	// The type of event, e.g. the topic
	Name *string `json:"name,omitempty"`
}

type InvokeDeploymentOperationRequest struct {
	Name *string `json:"name,omitempty"`

	// Anything
	Parameters any `json:"parameters,omitempty"`

	// whether the call wants to wait for the operation or if the result of the invocation will be received async from an event
// susbscription
	Wait *bool `json:"wait,omitempty"`
}

type InvokedDeploymentOperationResponse struct {
	DeploymentID *int32 `json:"deploymentId,omitempty"`
	ID *string `json:"id,omitempty"`
	InvokedOn *time.Time `json:"invokedOn,omitempty"`
	Name *string `json:"name,omitempty"`

	// Anything
	Parameters any `json:"parameters,omitempty"`

	// Anything
	Result any `json:"result,omitempty"`
	Status *string `json:"status,omitempty"`
}

// Operation - Defines an available operation
type Operation struct {
	Name *string `json:"name,omitempty"`
	Parameters []*OperationParameterType `json:"parameters,omitempty"`
	Target *OperationTargetType `json:"target,omitempty"`
}

// OperationParameterType - The parameter type information for a parameter of an operation
type OperationParameterType struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

type OperationTargetType struct {
	IDType *string `json:"idType,omitempty"`
	ObjectType *string `json:"objectType,omitempty"`
}

