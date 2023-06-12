package operation

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

type Scheduler interface {
	Schedule(ctx context.Context, operationId uuid.UUID) error
}

type scheduler struct {
	sender messaging.MessageSender
}

// Schedule schedules an operation for execution
func (scheduler *scheduler) Schedule(ctx context.Context, operationId uuid.UUID) error {
	message := ExecuteOperationCommand{OperationId: operationId}

	results, err := scheduler.sender.Send(ctx, string(messaging.QueueNameOperations), message)
	if err != nil {
		return err
	}

	if len(results) == 1 && results[0].Error != nil {
		return results[0].Error
	}
	return nil
}

func NewScheduler(sender messaging.MessageSender) (Scheduler, error) {
	return &scheduler{
		sender: sender,
	}, nil
}

func NewSchedulerFromConfig(appConfig *config.AppConfig) (Scheduler, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
		SubscriptionId:          appConfig.Azure.SubscriptionId,
		Location:                appConfig.Azure.Location,
		ResourceGroupName:       appConfig.Azure.ResourceGroupName,
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	})

	if err != nil {
		return nil, err
	}

	scheduler, err := NewScheduler(sender)
	if err != nil {
		return nil, err
	}
	return scheduler, nil
}
