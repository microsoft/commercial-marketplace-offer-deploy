package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func SetupResourceGroup(subscriptionId string, resourceGroupName string, location string) {
	resp, err := CreateResourceGroup(subscriptionId, resourceGroupName, location)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%s was created", *resp.Name)
	}
}

func DoesResourceGroupExist(subscriptionId string, resourceGroupName string, location string) (bool, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Print(err)
	}
	ctx := context.Background()

	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		return true, err
	}

	resp, err := resourceGroupClient.CheckExistence(ctx, resourceGroupName, nil)
	if err != nil {
		return true, err
	}

	return resp.Success, nil
}

func CreateResourceGroup(subscriptionId string, resourceGroupName string, location string) (*armresources.ResourceGroup, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Print(err)
	}
	ctx := context.Background()

	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(location),
		},
		nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}

func DeployPolicyDefinition(subscriptionId string) {
	log.Printf("Inside deployPolicyDefinition()")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armpolicy.NewDefinitionsClient(subscriptionId, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.CreateOrUpdate(ctx,
		"ResourceNaming",
		armpolicy.Definition{
			Properties: &armpolicy.DefinitionProperties{
				Description: to.Ptr("Force resource names to begin with given 'prefix' and/or end with given 'suffix'"),
				DisplayName: to.Ptr("Enforce resource naming convention"),
				Metadata: map[string]interface{}{
					"category": "Naming",
				},
				Mode: to.Ptr("All"),
				// Parameters: map[string]*armpolicy.ParameterDefinitionsValue{
				// 	"prefix": {
				// 		Type: to.Ptr(armpolicy.ParameterTypeString),
				// 		Metadata: &armpolicy.ParameterDefinitionsValueMetadata{
				// 			Description: to.Ptr("Resource name prefix"),
				// 			DisplayName: to.Ptr("Prefix"),
				// 		},
				// 	},
				// 	"suffix": {
				// 		Type: to.Ptr(armpolicy.ParameterTypeString),
				// 		Metadata: &armpolicy.ParameterDefinitionsValueMetadata{
				// 			Description: to.Ptr("Resource name suffix"),
				// 			DisplayName: to.Ptr("Suffix"),
				// 		},
				// 	},
				// },
				PolicyRule: map[string]interface{}{
					"if": map[string]interface{}{
						"not": map[string]interface{}{
							"field": "name",
							"like":  "a*b",
							//"like":  "[concat(parameters('prefix'), '*', parameters('suffix'))]",
						},
					},
					"then": map[string]interface{}{
						"effect": "deny",
					},
				},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func DeployPolicy(subscriptionId string, resourceGroupName string) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armpolicy.NewAssignmentsClient(subscriptionId, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, resourceGroupName)
	log.Printf("scope is %s", scope)

	policyDefinitionId := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/policyDefinitions/ResourceNaming", subscriptionId)
	log.Printf("policyDefinitionId is %s", policyDefinitionId)

	_, err = client.Create(ctx,
		scope,
		"ResourceName",
		armpolicy.Assignment{
			Properties: &armpolicy.AssignmentProperties{
				Description: to.Ptr("Enforce resource naming conventions"),
				DisplayName: to.Ptr("Enforce Resource Names"),
				Scope:       &scope,
				Metadata: map[string]interface{}{
					"assignedBy": "John Doe",
				},
				NonComplianceMessages: []*armpolicy.NonComplianceMessage{
					{
						Message: to.Ptr("A resource name was non-complaint.  It must be in the format 'a*b'."),
					}},
				PolicyDefinitionID: to.Ptr(policyDefinitionId),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
