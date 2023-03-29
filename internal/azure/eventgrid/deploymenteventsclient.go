package eventgrid

import (
	"context"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

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

func (c *deploymentEventsClient) CreateEventSubscription(ctx context.Context) error {
	_, err := armeventgrid.NewEventSubscriptionsClient(c.Properties.SubscriptionId, c.Credential, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *deploymentEventsClient) CreateSystemTopic(ctx context.Context) (*armeventgrid.SystemTopic, error) {
	systemTopicsClient, err := armeventgrid.NewSystemTopicsClient(c.Properties.SubscriptionId, c.Credential, nil)
	if err != nil {
		return nil, err
	}

	log.Print(c.Properties.Location)

	pollerResp, err := systemTopicsClient.BeginCreateOrUpdate(
		ctx,
		c.Properties.ResourceGroupName,
		c.Properties.SystemTopicName,
		armeventgrid.SystemTopic{
			Location: &c.Properties.Location,
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
	Location          string
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
	}

	resourceGroupClient, err := armresources.NewResourceGroupsClient(props.SubscriptionId, cred, nil)

	if err != nil {
		return nil, err
	}

	resourceGroup, err := resourceGroupClient.Get(ctx, props.ResourceGroupName, nil)

	if err != nil {
		return nil, err
	}

	props.Location = *resourceGroup.Location
	props.SystemTopicName = props.ResourceGroupName + "-event-topic"

	return props, nil
}
