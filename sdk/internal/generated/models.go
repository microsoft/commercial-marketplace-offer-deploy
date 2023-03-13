//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package generated

import "time"

type AvailableDeploymentOperation struct {
	Name *string `json:"name,omitempty"`
	Parameters []*string `json:"parameters,omitempty"`
}

type Deployment struct {
	ID *int64 `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Stages []*DeploymentStage `json:"stages,omitempty"`
	Template *DeploymentTemplate `json:"template,omitempty"`
}

// DeploymentManagementClientAddDeploymentOptions contains the optional parameters for the DeploymentManagementClient.AddDeployment
// method.
type DeploymentManagementClientAddDeploymentOptions struct {
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

// DeploymentManagementClientListOperationsOptions contains the optional parameters for the DeploymentManagementClient.ListOperations
// method.
type DeploymentManagementClientListOperationsOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientRegisterEventSubscriptionOptions contains the optional parameters for the DeploymentManagementClient.RegisterEventSubscription
// method.
type DeploymentManagementClientRegisterEventSubscriptionOptions struct {
	// placeholder for future optional parameters
}

// DeploymentManagementClientUpdateDeploymentOptions contains the optional parameters for the DeploymentManagementClient.UpdateDeployment
// method.
type DeploymentManagementClientUpdateDeploymentOptions struct {
	// placeholder for future optional parameters
}

type DeploymentOperation struct {
	ID *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Parameters []*DeploymentOperationParameter `json:"parameters,omitempty"`
	TriggeredOn *time.Time `json:"triggeredOn,omitempty"`
}

type DeploymentOperationParameter struct {
	Name *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

type DeploymentRequest struct {
	// REQUIRED
	Name *string `json:"name,omitempty"`
	MultiStage *bool `json:"multiStage,omitempty"`
	Template *DeploymentTemplate `json:"template,omitempty"`
}

type DeploymentStage struct {
	// Dictionary of
	Attributes map[string]any `json:"attributes,omitempty"`
	Name *string `json:"name,omitempty"`
}

type DeploymentTemplate struct {
	Name *string `json:"name,omitempty"`
	Parameters []*DeploymentTemplateParameter `json:"parameters,omitempty"`
	URI *string `json:"uri,omitempty"`
}

type DeploymentTemplateParameter struct {
	Name *string `json:"name,omitempty"`

	// Dictionary of
	Value map[string]any `json:"value,omitempty"`
}

type Event struct {
	Name *string `json:"name,omitempty"`
}

type EventSubscription struct {
	Callback *string `json:"callback,omitempty"`
	Name *string `json:"name,omitempty"`
	Topic *string `json:"topic,omitempty"`
}

type InvokeDeploymentOperation struct {
	Name *string `json:"name,omitempty"`
	Parameters []*DeploymentOperationParameter `json:"parameters,omitempty"`
}

type InvokeDeploymentOperationResult struct {
	ID *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

