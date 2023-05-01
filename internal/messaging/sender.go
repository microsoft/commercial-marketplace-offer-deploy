package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
)

type MessageSenderOptions struct {
	SubscriptionId          string
	Location                string
	ResourceGroupName       string
	FullyQualifiedNamespace string
}

func (o *MessageSenderOptions) getNamespaceName() string {
	parts := strings.Split(o.FullyQualifiedNamespace, ".")
	return parts[0]
}

type SendMessageResult struct {
	Success bool
	Error   error
}

type MessageSender interface {
	Send(ctx context.Context, queueName string, messages ...any) ([]SendMessageResult, error)
	EnsureTopology(ctx context.Context, queueName string) error
}

type serviceBusMessageSender struct {
	options       MessageSenderOptions
	client        *azservicebus.Client
	clientFactory *armservicebus.ClientFactory
	topology      bool
}

// TODO: need to get the message bus from configuration

func NewServiceBusMessageSender(credential azcore.TokenCredential, options MessageSenderOptions) (MessageSender, error) {
	log.Debugf("New Service Bus Message Sender options:\n %v", options)

	client, err := azservicebus.NewClient(options.FullyQualifiedNamespace, credential, nil)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armservicebus.NewClientFactory(options.SubscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	return &serviceBusMessageSender{options: options, client: client, clientFactory: clientFactory}, nil
}

func (s *serviceBusMessageSender) Send(ctx context.Context, queueName string, messages ...any) ([]SendMessageResult, error) {
	sender, err := s.client.NewSender(queueName, nil)

	if err != nil {
		return nil, err
	}
	defer sender.Close(ctx)

	results := []SendMessageResult{}

	for index, message := range messages {
		body, err := json.Marshal(message)
		log.Debugf("marshaling message %d:\n %s", index, string(body))

		if err != nil {
			results = append(results, SendMessageResult{Success: false, Error: fmt.Errorf("failed to marshal message %d: %w", index, err)})
			continue
		}

		err = sender.SendMessage(ctx, &azservicebus.Message{
			Body: body,
		}, nil)
		log.Debug("sent message")

		results = append(results, SendMessageResult{Success: err == nil, Error: err})
	}
	log.Debugf("finished sending messages with results of %v", results)
	return results, nil
}

// implement a method on serviceBusMessageSender that ensures the queue exists
func (s *serviceBusMessageSender) EnsureTopology(ctx context.Context, queueName string) error {
	if s.topology {
		return nil
	}

	_, err := s.createOrUpdateNamespace(ctx)
	if err != nil {
		return err
	}

	_, err = s.createOrUpdateQueue(ctx, queueName)
	if err != nil {
		return err
	}
	s.topology = true
	return nil
}

func (s *serviceBusMessageSender) createOrUpdateNamespace(ctx context.Context) (*armservicebus.SBNamespace, error) {
	namespacesClient := s.clientFactory.NewNamespacesClient()

	pollerResp, err := namespacesClient.BeginCreateOrUpdate(
		ctx,
		s.options.ResourceGroupName,
		s.options.getNamespaceName(),
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
		s.options.getNamespaceName(),
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
