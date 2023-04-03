package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
)

type MessageSenderOptions struct {
	SubscriptionId      string
	Location            string
	ResourceGroupName   string
	ServiceBusNamespace string
}

type SendMessageResult struct {
	Success bool
	Error   error
}

type MessageSender interface {
	Send(ctx context.Context, queueName string, messages ...any) ([]SendMessageResult, error)
}

type serviceBusMessageSender struct {
	options       MessageSenderOptions
	client        *azservicebus.Client
	clientFactory *armservicebus.ClientFactory
}

// TODO: need to get the message bus from configuration

func NewServiceBusMessageSender(options MessageSenderOptions, credential azcore.TokenCredential) (MessageSender, error) {
	client, err := azservicebus.NewClient(options.ServiceBusNamespace, credential, nil)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armservicebus.NewClientFactory(options.SubscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	return &serviceBusMessageSender{client: client, clientFactory: clientFactory}, nil
}

func (s *serviceBusMessageSender) Send(ctx context.Context, queueName string, messages ...any) ([]SendMessageResult, error) {
	err := s.ensureTopology(ctx, queueName)
	if err != nil {
		return nil, err
	}

	sender, err := s.client.NewSender(queueName, nil)

	if err != nil {
		return nil, err
	}
	defer sender.Close(ctx)

	results := []SendMessageResult{}

	for index, message := range messages {
		body, err := json.Marshal(message)
		if err != nil {
			results = append(results, SendMessageResult{Success: false, Error: fmt.Errorf("failed to marshal message %d: %w", index, err)})
			continue
		}

		err = sender.SendMessage(ctx, &azservicebus.Message{
			Body: body,
		}, nil)

		results = append(results, SendMessageResult{Success: err == nil, Error: err})
	}
	return results, nil
}

// implement a method on serviceBusMessageSender that ensures the queue exists
func (s *serviceBusMessageSender) ensureTopology(ctx context.Context, queueName string) error {
	_, err := s.createOrUpdateNamespace(ctx)
	if err != nil {
		return err
	}

	_, err = s.createOrUpdateQueue(ctx, queueName)
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceBusMessageSender) createOrUpdateNamespace(ctx context.Context) (*armservicebus.SBNamespace, error) {
	namespacesClient := s.clientFactory.NewNamespacesClient()

	pollerResp, err := namespacesClient.BeginCreateOrUpdate(
		ctx,
		s.options.ResourceGroupName,
		s.options.ServiceBusNamespace,
		armservicebus.SBNamespace{
			Location: to.Ptr(s.options.Location),
			SKU: &armservicebus.SBSKU{
				Name: to.Ptr(armservicebus.SKUNameStandard),
				Tier: to.Ptr(armservicebus.SKUTierStandard),
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
	return &resp.SBNamespace, nil
}

func (s *serviceBusMessageSender) createOrUpdateQueue(ctx context.Context, queueName string) (*armservicebus.SBQueue, error) {
	queuesClient := s.clientFactory.NewQueuesClient()

	resp, err := queuesClient.CreateOrUpdate(
		ctx,
		s.options.ResourceGroupName,
		s.options.ServiceBusNamespace,
		queueName,
		armservicebus.SBQueue{
			Properties: &armservicebus.SBQueueProperties{
				EnablePartitioning:                  to.Ptr(false),
				RequiresDuplicateDetection:          to.Ptr(true),
				DuplicateDetectionHistoryTimeWindow: to.Ptr("PT10M"),
			},
		},
		nil,
	)
	if err != nil {
		return nil, nil
	}

	return &resp.SBQueue, nil
}
