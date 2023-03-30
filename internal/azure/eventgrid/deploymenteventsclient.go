package eventgrid

import (
	"context"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid"
)

// System topics MUST be set to global location
const systemTopicLocation = "global"

type DeploymentEventsClient interface {
	CreateSystemTopic(ctx context.Context) (*armeventgrid.SystemTopic, error)
	CreateEventSubscription(ctx context.Context) error
}

type deploymentEventsClient struct {
	Credential azcore.TokenCredential
	Properties *eventGridManagerProperties
}

// Creates an event grid manager to create system topic and event subscription for the purpose of receiving deployment events
// It will use the resource group id to create the system topic in the resource group's location and subscription
func NewDeploymentEventsClient(credential azcore.TokenCredential, resourceGroupId string) (DeploymentEventsClient, error) {
	properties, err := getProperties(context.TODO(), credential, resourceGroupId)

	if err != nil {
		return nil, err
	}
	client := &deploymentEventsClient{
		Credential: credential,
		Properties: properties,
	}
	return client, nil
}

func (c *deploymentEventsClient) CreateEventSubscription(ctx context.Context) (*armeventgrid.SystemTopicEventSubscriptionsClientGetResponse, error) {
	// TODO: change '_' to eventSubscriptionsClient
	// next: using the event topic from Properties, create an event subscription (web hook) using the endpoint of the operator
	// parameters required:
	//	endpointUrl string: the endpoint of the operator should be a parameter on method CreateEventSubscription
	//  subscriptionName string: this is going to be the name of the subscription
	// example of an event subscription for a system topic:
	//		https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/resourcemanager/eventgrid/armeventgrid/ze_generated_example_systemtopiceventsubscriptions_client_test.go

	// return: modify the return from error to a tuple (*armeventgrid.[whatever the response is from the Azure client], error)

	// Step:2
	// inside this file, create a global package variable called "includedEventTypes" that has all the event types that we're interested
	// slice, populated with the below represented values:
	/*
	   "filter": {
	       "includedEventTypes": [
	           "Microsoft.Resources.ResourceWriteSuccess",
	           "Microsoft.Resources.ResourceWriteFailure",
	           "Microsoft.Resources.ResourceWriteCancel",
	           "Microsoft.Resources.ResourceDeleteSuccess",
	           "Microsoft.Resources.ResourceDeleteFailure",
	           "Microsoft.Resources.ResourceDeleteCancel",
	           "Microsoft.Resources.ResourceActionSuccess",
	           "Microsoft.Resources.ResourceActionFailure",
	           "Microsoft.Resources.ResourceActionCancel"
	       ],
	       "enableAdvancedFilteringOnArrays": true
	*/
	ctx = context.Background()
	eventSubscriptionsClient, err := armeventgrid.NewEventSubscriptionsClient(c.Properties.SubscriptionId, c.Credential, nil)
	if err != nil {
		return nil, err
	}

	pollerResp, err := eventSubscriptionsClient.BeginCreateOrUpdate(ctx, c.Properties.SubscriptionId, "ashwinse-test", armeventgrid.EventSubscription{
		Properties: &armeventgrid.EventSubscriptionProperties{
			Destination: &armeventgrid.WebHookEventSubscriptionDestination{
				EndpointType: to.Ptr(armeventgrid.EndpointTypeWebHook),
				Properties: &armeventgrid.WebHookEventSubscriptionDestinationProperties{
					EndpointURL: to.Ptr(" https://7f53-2601-600-8d81-1470-b9a0-ca77-9e5c-b6b5.ngrok.io/eventgrid"),
				},
			},
		},
	}, nil)

	if err != nil {
		log.Fatalf("failed to finish request: %v", err)
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}

	return &armeventgrid.SystemTopicEventSubscriptionsClientGetResponse{}, nil
	// _, err := armeventgrid.NewEventSubscriptionsClient(c.Properties.SubscriptionId, c.Credential, nil)
	// if err != nil {
	// 	return err
	// }
	// return nil
}

func (c *deploymentEventsClient) DeleteSystemTopic(ctx context.Context) (*armeventgrid.SystemTopicsClientDeleteResponse, error) {
	systemTopicsClient, err := armeventgrid.NewSystemTopicsClient(c.Properties.SubscriptionId, c.Credential, nil)
	if err != nil {
		return nil, err
	}

	pollerResp, err := systemTopicsClient.BeginDelete(
		ctx,
		c.Properties.ResourceGroupName,
		c.Properties.SystemTopicName,
		nil)

	if err != nil {
		return nil, err
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *deploymentEventsClient) CreateSystemTopic(ctx context.Context) (*armeventgrid.SystemTopic, error) {
	systemTopicsClient, err := armeventgrid.NewSystemTopicsClient(c.Properties.SubscriptionId, c.Credential, nil)
	if err != nil {
		return nil, err
	}

	pollerResp, err := systemTopicsClient.BeginCreateOrUpdate(
		ctx,
		c.Properties.ResourceGroupName,
		c.Properties.SystemTopicName,
		armeventgrid.SystemTopic{
			Location: to.Ptr(systemTopicLocation),
			Properties: &armeventgrid.SystemTopicProperties{
				Source:    to.Ptr(c.Properties.ResourceGroupId),
				TopicType: to.Ptr("Microsoft.Resources.ResourceGroups"),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.SystemTopic, nil
}

type eventGridManagerProperties struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceGroupId   string
	SystemTopicName   string
}

func getProperties(ctx context.Context, cred azcore.TokenCredential, resourceGroupId string) (*eventGridManagerProperties, error) {
	values := strings.Split(resourceGroupId, "/")
	props := &eventGridManagerProperties{
		SubscriptionId:    values[2],
		ResourceGroupName: values[len(values)-1],
		ResourceGroupId:   resourceGroupId,
		SystemTopicName:   values[len(values)-1] + "-event-topic",
	}

	return props, nil
}
