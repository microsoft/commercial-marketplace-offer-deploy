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
	options       *ServiceBusHealthCheckOptions
	sendResult    HealthCheckResult
	receiveResult HealthCheckResult
}

// Check whether the url is accessible
func (c *serviceBusHealthCheck) Check(ctx context.Context) HealthCheckResult {
	threshold := time.Now().Add(c.options.Timeout)

	for {
		if time.Now().After(threshold) {
			result := HealthCheckResult{
				Description: fmt.Sprintf("Timeout exceeded while connecting to service bus %s", c.options.FullyQualifiedNamespace),
				Status:      HealthCheckStatusUnhealthy,
				Error:       errors.New("timeout exceeded while waiting for connectivity to service"),
			}
			log.Warnf("Health Check timed out: %v", result)
			return result
		}
		result := c.getResult(ctx)

		if result.Status != HealthCheckStatusHealthy || result.Error != nil {
			log.Warnf("Health Check attempt failed: %v", result)
			time.Sleep(10 * time.Second)
			continue
		}

		return result
	}
}

func (c *serviceBusHealthCheck) getResult(ctx context.Context) HealthCheckResult {
	err := c.checkSend(ctx)

	if err != nil {
		return HealthCheckResult{
			Description: fmt.Sprintf("Failed to send message to queue %s", c.options.QueueName),
			Status:      HealthCheckStatusUnhealthy,
			Error:       err,
		}
	}

	err = c.checkReceiver(ctx)

	if err != nil {
		return HealthCheckResult{
			Description: fmt.Sprintf("Failed to receive message from queue %s", c.options.QueueName),
			Status:      HealthCheckStatusUnhealthy,
			Error:       err,
		}
	}

	if c.sendResult.Status == HealthCheckStatusHealthy && c.receiveResult.Status == HealthCheckStatusHealthy {
		return HealthCheckResult{
			Description: fmt.Sprintf("Successfully connected to queue %s", c.options.QueueName),
			Status:      HealthCheckStatusHealthy,
		}
	}

	return HealthCheckResult{
		Description: fmt.Sprintf("Failed to connect to queue %s", c.options.QueueName),
		Status:      HealthCheckStatusUnhealthy,
		Error:       errors.New("failed to connect to queue"),
	}
}

func (c *serviceBusHealthCheck) checkSend(ctx context.Context) error {
	if c.sendResult.Status == HealthCheckStatusHealthy {
		return nil
	}

	client, err := c.getServiceBusClient(ctx)
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

	if err == nil {
		c.sendResult = HealthCheckResult{
			Description: fmt.Sprintf("Successfully sent message to queue %s", c.options.QueueName),
			Status:      HealthCheckStatusHealthy,
		}
		log.Info(c.sendResult.Description)
	}

	return err
}

func (c *serviceBusHealthCheck) checkReceiver(ctx context.Context) error {
	if c.receiveResult.Status == HealthCheckStatusHealthy {
		return nil
	}

	client, err := c.getServiceBusClient(ctx)
	if err != nil {
		return err
	}

	receiver, err := client.NewReceiverForQueue(c.options.QueueName, nil)
	if err != nil {
		return err
	}
	defer receiver.Close(ctx)

	_, err = receiver.ReceiveMessages(ctx, 100, nil)

	if err != nil {
		return err
	} else {
		c.receiveResult = HealthCheckResult{
			Description: fmt.Sprintf("Successfully received message from queue %s", c.options.QueueName),
			Status:      HealthCheckStatusHealthy,
		}
		log.Info(c.receiveResult.Description)
	}

	return nil
}

func (c *serviceBusHealthCheck) getServiceBusClient(ctx context.Context) (*azservicebus.Client, error) {
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
