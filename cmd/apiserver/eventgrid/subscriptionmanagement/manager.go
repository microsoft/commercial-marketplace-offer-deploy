package subscriptionmanagement

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid"
)

// System topics MUST be set to global location
const systemTopicLocation = "global"

type EventGridManager interface {
	GetSystemTopicName() string
	CreateSystemTopic(ctx context.Context) (*armeventgrid.SystemTopic, error)
	CreateEventSubscription(ctx context.Context, subscriptionName string, endpointUrl string) (*armeventgrid.EventSubscriptionsClientCreateOrUpdateResponse, error)
}

// internal implementation of event grid manager
type manager struct {
	Credential azcore.TokenCredential
	Properties *eventGridManagerProperties
}

func (c *manager) GetSystemTopicName() string {
	return c.Properties.SystemTopicName
}

// Creates an event grid manager to create system topic and event subscription for the purpose of receiving deployment events
// It will use the resource group id to create the system topic in the resource group's location and subscription
func NewEventGridManager(credential azcore.TokenCredential, resourceGroupId string) (EventGridManager, error) {
	properties, err := getProperties(context.TODO(), credential, resourceGroupId)
	if err != nil {
		return nil, err
	}

	client := &manager{
		Credential: credential,
		Properties: properties,
	}
	return client, nil
}

func (c *manager) CreateEventSubscription(ctx context.Context, subscriptionName string, endpointUrl string) (*armeventgrid.EventSubscriptionsClientCreateOrUpdateResponse, error) {
	eventSubscriptionsClient, err := armeventgrid.NewEventSubscriptionsClient(c.Properties.SubscriptionId, c.Credential, nil)
	if err != nil {
		return nil, err
	}

	subscriptionScope := c.Properties.ResourceGroupId
	filter := getDeploymentResourceSubscriptionFilter()

	properties := &armeventgrid.EventSubscriptionProperties{
		Destination: &armeventgrid.WebHookEventSubscriptionDestination{
			EndpointType: to.Ptr(armeventgrid.EndpointTypeWebHook),
			Properties: &armeventgrid.WebHookEventSubscriptionDestinationProperties{
				EndpointURL: to.Ptr(endpointUrl),
			},
		},
		Filter: filter,
	}

	pollerResp, err := eventSubscriptionsClient.BeginCreateOrUpdate(ctx,
		subscriptionScope,
		subscriptionName,
		armeventgrid.EventSubscription{Properties: properties}, nil)

	if err != nil {
		return nil, err
	}

	response, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

func (c *manager) DeleteSystemTopic(ctx context.Context) (*armeventgrid.SystemTopicsClientDeleteResponse, error) {
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

func (c *manager) CreateSystemTopic(ctx context.Context) (*armeventgrid.SystemTopic, error) {
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
		if responseError, ok := err.(*azcore.ResponseError); ok {
			if responseError.StatusCode == 400 && strings.Contains(err.Error(), "Only one system topic is allowed per source.") {
				log.Print("System topic already exists for resource group")
				return nil, nil
			}
		} else {
			return nil, err
		}
	}

	if pollerResp != nil {
		resp, err := pollerResp.PollUntilDone(ctx, nil)
		if err != nil {
			return nil, err
		}

		log.Printf("Created system topic %s in resource group %s", c.Properties.SystemTopicName, c.Properties.ResourceGroupName)
		return &resp.SystemTopic, nil
	}
	return nil, fmt.Errorf("poller response nil in Event Grid subscription manager")
}

type eventGridManagerProperties struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceGroupId   string
	SystemTopicName   string
}

func getDeploymentResourceSubscriptionFilter() *armeventgrid.EventSubscriptionFilter {
	return &armeventgrid.EventSubscriptionFilter{
		EnableAdvancedFilteringOnArrays: to.Ptr(true),
		IncludedEventTypes:              getIncludedEventTypesForFilter(),
	}
}

func getIncludedEventTypesForFilter() []*string {
	return []*string{
		// filter on what we care about (what a consumer can take action on)
		// we don't need to worry about message ordering if we only listen for success and failure

		// if this gets modified for more than success and failure, message ordering will need to be considered
		to.Ptr("Microsoft.Resources.ResourceWriteSuccess"),
		to.Ptr("Microsoft.Resources.ResourceWriteFailure"),
		to.Ptr("Microsoft.Resources.ResourceDeleteSuccess"),
		to.Ptr("Microsoft.Resources.ResourceDeleteFailure"),
	}
}

func getProperties(ctx context.Context, cred azcore.TokenCredential, resourceGroupId string) (*eventGridManagerProperties, error) {
	resourceId, err := arm.ParseResourceID(resourceGroupId)

	if err != nil {
		return nil, err
	}
	props := &eventGridManagerProperties{
		SubscriptionId:    resourceId.SubscriptionID,
		ResourceGroupName: resourceId.ResourceGroupName,
		ResourceGroupId:   resourceGroupId,
		SystemTopicName:   resourceId.ResourceGroupName + "-event-topic",
	}
	return props, nil
}
