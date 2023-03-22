//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package internal

import "time"

type CreateDeployment struct {
	// REQUIRED
	Name *string `json:"name,omitempty"`

	// REQUIRED; Anything
	Template any `json:"template,omitempty"`
}

type CreateEventSubscription struct {
	Callback *string `json:"callback,omitempty"`
	Name *string `json:"name,omitempty"`
}

type Deployment struct {
	ID *int32 `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
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

// DeploymentManagementClientGetDeploymentOperationOptions contains the optional parameters for the DeploymentManagementClient.GetDeploymentOperation
// method.
type DeploymentManagementClientGetDeploymentOperationOptions struct {
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

// DeploymentManagementClientGetEventsOptions contains the optional parameters for the DeploymentManagementClient.GetEvents
// method.
type DeploymentManagementClientGetEventsOptions struct {
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

type Event struct {
	Name *string `json:"name,omitempty"`
}

type EventSubscription struct {
	Callback *string `json:"callback,omitempty"`
	ID *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Topic *string `json:"topic,omitempty"`
}

type InvokeDeploymentOperation struct {
	Name *string `json:"name,omitempty"`
	
	//Parameters map[string] interface{} `json:"parameters,omitempty"`

	// Anything
	Parameters any `json:"parameters,omitempty"`

	// whether the call wants to wait for the operation or if the result of the invocation will be received async from an event
// susbscription
	Wait *bool `json:"wait,omitempty"`
}

type InvokedOperation struct {
	ID *string `json:"id,omitempty"`
	InvokedOn *time.Time `json:"invokedOn,omitempty"`
	Name *string `json:"name,omitempty"`

	// Anything
	Parameters any `json:"parameters,omitempty"`

	// Anything
	Result any `json:"result,omitempty"`
	Status *string `json:"status,omitempty"`
	Target *InvokedOperationTarget `json:"target,omitempty"`
}

type InvokedOperationTarget struct {
	ID *InvokedOperationTargetID `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
}

type InvokedOperationTargetID struct {
	Type *string `json:"type,omitempty"`
	Value *string `json:"value,omitempty"`
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

