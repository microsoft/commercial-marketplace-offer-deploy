package diagnostics

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	log "github.com/sirupsen/logrus"
)

const HealthCheckQueueName = "healthcheck"

type ServiceBusHealthCheckOptions struct {
	Timeout                 time.Duration
	QueueName               string
	FullyQualifiedNamespace string
}

type serviceBusHealthCheck struct {
	options *ServiceBusHealthCheckOptions
}

// Check whether the url is accessible
func (c *serviceBusHealthCheck) Check(ctx context.Context) HealthCheckResult {
	threshold := time.Now().Add(c.options.Timeout)

	for {
		if time.Now().After(threshold) {
			return HealthCheckResult{
				Description: fmt.Sprintf("Timeout exceeded while connecting to service bus %s", c.options.FullyQualifiedNamespace),
				Status:      HealthCheckStatusUnhealthy,
				Error:       errors.New("timeout exceeded while waiting for connectivity to service"),
			}
		}
		result := c.getResult(ctx)

		if result.Status != HealthCheckStatusHealthy || result.Error != nil {
			log.Warnf("Health Check attempt failed: %v", result)
			time.Sleep(5 * time.Second)
			continue
		}

		return result
	}
}

func (c *serviceBusHealthCheck) getResult(ctx context.Context) HealthCheckResult {
	err := c.checkSend()

	if err != nil {
		return HealthCheckResult{
			Description: fmt.Sprintf("Failed to send message to queue %s", c.options.QueueName),
			Status:      HealthCheckStatusUnhealthy,
			Error:       err,
		}
	}

	err = c.checkReceiver()

	if err != nil {
		return HealthCheckResult{
			Description: fmt.Sprintf("Failed to receive message from queue %s", c.options.QueueName),
			Status:      HealthCheckStatusUnhealthy,
			Error:       err,
		}
	}

	return HealthCheckResult{
		Description: fmt.Sprintf("Successfully sent and received message from queue %s", c.options.QueueName),
		Status:      HealthCheckStatusHealthy,
	}
}

func (c *serviceBusHealthCheck) checkSend() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := c.getServiceBusClient()
	if err != nil {
		return err
	}

	sender, err := client.NewSender(c.options.QueueName, nil)
	if err != nil {
		return err
	}
	defer sender.Close(ctx)

	err = sender.SendMessage(ctx, &azservicebus.Message{
		Body: []byte("test"),
	}, nil)

	return err
}

func (c *serviceBusHealthCheck) checkReceiver() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := c.getServiceBusClient()
	if err != nil {
		return err
	}

	receiver, err := client.NewReceiverForQueue(c.options.QueueName, nil)
	if err != nil {
		return err
	}
	defer receiver.Close(ctx)

	var messages []*azservicebus.ReceivedMessage
	messages, err = receiver.ReceiveMessages(ctx, 1, nil)

	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return errors.New("no messages received")
	}

	return nil
}

func (c *serviceBusHealthCheck) getServiceBusClient() (*azservicebus.Client, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, err
	}

	return azservicebus.NewClient(c.options.FullyQualifiedNamespace, credential, nil)
}

func NewServiceBusHealthCheck(options ServiceBusHealthCheckOptions) HealthCheck {

	return &serviceBusHealthCheck{
		options: &options,
	}
}
