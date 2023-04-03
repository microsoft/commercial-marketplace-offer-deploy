package subscriptionmanagement

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid"
)

// System topics MUST be set to global location
const systemTopicLocation = "global"

type EventGridManager interface {
	CreateSystemTopic(ctx context.Context) (*armeventgrid.SystemTopic, error)
	CreateEventSubscription(ctx context.Context, subscriptionName string, endpointUrl string) (*armeventgrid.EventSubscriptionsClientCreateOrUpdateResponse, error)
}

// internal implementation of event grid manager
type manager struct {
	Credential azcore.TokenCredential
	Properties *eventGridManagerProperties
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

func getDeploymentResourceSubscriptionFilter() *armeventgrid.EventSubscriptionFilter {
	return &armeventgrid.EventSubscriptionFilter{
		EnableAdvancedFilteringOnArrays: to.Ptr(true),
		IncludedEventTypes:              getIncludedEventTypesForFilter(),
		AdvancedFilters: []armeventgrid.AdvancedFilterClassification{
			&armeventgrid.StringContainsAdvancedFilter{
				Values: []*string{
					to.Ptr("/providers/Microsoft.Resources/deployments/"),
				},
				OperatorType: to.Ptr(armeventgrid.AdvancedFilterOperatorTypeStringContains),
				Key:          to.Ptr("subject"),
			},
		},
	}
}

func getIncludedEventTypesForFilter() []*string {
	return []*string{
		to.Ptr("Microsoft.Resources.ResourceWriteSuccess"),
		to.Ptr("Microsoft.Resources.ResourceWriteFailure"),
		to.Ptr("Microsoft.Resources.ResourceWriteCancel"),
		to.Ptr("Microsoft.Resources.ResourceDeleteSuccess"),
		to.Ptr("Microsoft.Resources.ResourceDeleteFailure"),
		to.Ptr("Microsoft.Resources.ResourceDeleteCancel"),
		to.Ptr("Microsoft.Resources.ResourceActionSuccess"),
		to.Ptr("Microsoft.Resources.ResourceActionFailure"),
		to.Ptr("Microsoft.Resources.ResourceActionCancel"),
	}
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
