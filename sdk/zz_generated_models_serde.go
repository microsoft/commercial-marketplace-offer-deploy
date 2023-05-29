//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type CreateDeployment.
func (c CreateDeployment) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "location", c.Location)
	populate(objectMap, "name", c.Name)
	populate(objectMap, "resourceGroup", c.ResourceGroup)
	populate(objectMap, "subscriptionId", c.SubscriptionID)
	populate(objectMap, "template", &c.Template)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type CreateDeployment.
func (c *CreateDeployment) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "location":
				err = unpopulate(val, "Location", &c.Location)
				delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &c.Name)
				delete(rawMsg, key)
		case "resourceGroup":
				err = unpopulate(val, "ResourceGroup", &c.ResourceGroup)
				delete(rawMsg, key)
		case "subscriptionId":
				err = unpopulate(val, "SubscriptionID", &c.SubscriptionID)
				delete(rawMsg, key)
		case "template":
				err = unpopulate(val, "Template", &c.Template)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type CreateEventHookRequest.
func (c CreateEventHookRequest) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "apiKey", c.APIKey)
	populate(objectMap, "callback", c.Callback)
	populate(objectMap, "name", c.Name)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type CreateEventHookRequest.
func (c *CreateEventHookRequest) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "apiKey":
				err = unpopulate(val, "APIKey", &c.APIKey)
				delete(rawMsg, key)
		case "callback":
				err = unpopulate(val, "Callback", &c.Callback)
				delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &c.Name)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type CreateEventHookResponse.
func (c CreateEventHookResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "id", c.ID)
	populate(objectMap, "name", c.Name)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type CreateEventHookResponse.
func (c *CreateEventHookResponse) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "id":
				err = unpopulate(val, "ID", &c.ID)
				delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &c.Name)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type Deployment.
func (d Deployment) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "id", d.ID)
	populate(objectMap, "name", d.Name)
	populate(objectMap, "stages", d.Stages)
	populate(objectMap, "status", d.Status)
	populate(objectMap, "template", &d.Template)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Deployment.
func (d *Deployment) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", d, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "id":
				err = unpopulate(val, "ID", &d.ID)
				delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &d.Name)
				delete(rawMsg, key)
		case "stages":
				err = unpopulate(val, "Stages", &d.Stages)
				delete(rawMsg, key)
		case "status":
				err = unpopulate(val, "Status", &d.Status)
				delete(rawMsg, key)
		case "template":
				err = unpopulate(val, "Template", &d.Template)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", d, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type DeploymentStage.
func (d DeploymentStage) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "deploymentName", d.DeploymentName)
	populate(objectMap, "id", d.ID)
	populate(objectMap, "name", d.Name)
	populate(objectMap, "retries", d.Retries)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DeploymentStage.
func (d *DeploymentStage) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", d, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "deploymentName":
				err = unpopulate(val, "DeploymentName", &d.DeploymentName)
				delete(rawMsg, key)
		case "id":
				err = unpopulate(val, "ID", &d.ID)
				delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &d.Name)
				delete(rawMsg, key)
		case "retries":
				err = unpopulate(val, "Retries", &d.Retries)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", d, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type EventHookResponse.
func (e EventHookResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "callback", e.Callback)
	populate(objectMap, "id", e.ID)
	populate(objectMap, "name", e.Name)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type EventHookResponse.
func (e *EventHookResponse) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "callback":
				err = unpopulate(val, "Callback", &e.Callback)
				delete(rawMsg, key)
		case "id":
				err = unpopulate(val, "ID", &e.ID)
				delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &e.Name)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type EventType.
func (e EventType) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "name", e.Name)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type EventType.
func (e *EventType) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "name":
				err = unpopulate(val, "Name", &e.Name)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type GetInvokedOperationResponse.
func (g GetInvokedOperationResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "invokedOperation", g.InvokedOperation)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type GetInvokedOperationResponse.
func (g *GetInvokedOperationResponse) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", g, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "invokedOperation":
				err = unpopulate(val, "InvokedOperation", &g.InvokedOperation)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", g, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type InvokeDeploymentOperationRequest.
func (i InvokeDeploymentOperationRequest) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "name", i.Name)
	populate(objectMap, "parameters", &i.Parameters)
	populate(objectMap, "retries", i.Retries)
	populate(objectMap, "wait", i.Wait)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type InvokeDeploymentOperationRequest.
func (i *InvokeDeploymentOperationRequest) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", i, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "name":
				err = unpopulate(val, "Name", &i.Name)
				delete(rawMsg, key)
		case "parameters":
				err = unpopulate(val, "Parameters", &i.Parameters)
				delete(rawMsg, key)
		case "retries":
				err = unpopulate(val, "Retries", &i.Retries)
				delete(rawMsg, key)
		case "wait":
				err = unpopulate(val, "Wait", &i.Wait)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", i, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type InvokedDeploymentOperationResponse.
func (i InvokedDeploymentOperationResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "invokedOperation", i.InvokedOperation)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type InvokedDeploymentOperationResponse.
func (i *InvokedDeploymentOperationResponse) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", i, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "invokedOperation":
				err = unpopulate(val, "InvokedOperation", &i.InvokedOperation)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", i, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type InvokedOperation.
func (i InvokedOperation) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "attempts", i.Attempts)
	populate(objectMap, "deploymentId", i.DeploymentID)
	populate(objectMap, "id", i.ID)
	populateTimeRFC3339(objectMap, "invokedOn", i.InvokedOn)
	populate(objectMap, "name", i.Name)
	populate(objectMap, "parameters", &i.Parameters)
	populate(objectMap, "result", &i.Result)
	populate(objectMap, "retries", i.Retries)
	populate(objectMap, "status", i.Status)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type InvokedOperation.
func (i *InvokedOperation) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", i, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "attempts":
				err = unpopulate(val, "Attempts", &i.Attempts)
				delete(rawMsg, key)
		case "deploymentId":
				err = unpopulate(val, "DeploymentID", &i.DeploymentID)
				delete(rawMsg, key)
		case "id":
				err = unpopulate(val, "ID", &i.ID)
				delete(rawMsg, key)
		case "invokedOn":
				err = unpopulateTimeRFC3339(val, "InvokedOn", &i.InvokedOn)
				delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &i.Name)
				delete(rawMsg, key)
		case "parameters":
				err = unpopulate(val, "Parameters", &i.Parameters)
				delete(rawMsg, key)
		case "result":
				err = unpopulate(val, "Result", &i.Result)
				delete(rawMsg, key)
		case "retries":
				err = unpopulate(val, "Retries", &i.Retries)
				delete(rawMsg, key)
		case "status":
				err = unpopulate(val, "Status", &i.Status)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", i, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ListInvokedOperationResponse.
func (l ListInvokedOperationResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "items", l.Items)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ListInvokedOperationResponse.
func (l *ListInvokedOperationResponse) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", l, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "items":
				err = unpopulate(val, "Items", &l.Items)
				delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", l, err)
		}
	}
	return nil
}

func populate(m map[string]any, k string, v any) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}

func unpopulate(data json.RawMessage, fn string, v any) error {
	if data == nil {
		return nil
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}
	return nil
}
