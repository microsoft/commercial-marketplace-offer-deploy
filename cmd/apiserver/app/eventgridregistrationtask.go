package app

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/subscriptionmanagement"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
	log "github.com/sirupsen/logrus"
)

// constructor for creating task that registers event grid system topic for the resource group deployment events
func newEventGridRegistrationTask(appConfig *config.AppConfig) tasks.Task {
	taskOptions := eventGridRegistrationTaskOptions{
		Credential:      hosting.GetAzureCredential(),
		ResourceGroupId: appConfig.Azure.GetResourceGroupId(),
		EndpointUrl:     appConfig.GetPublicBaseUrl() + "eventgrid",
	}
	task := create(taskOptions)

	return task
}

type eventGridRegistrationTaskOptions struct {
	Credential      azcore.TokenCredential
	ResourceGroupId string
	EndpointUrl     string
	IsReady         func() bool
}

// factory for creating task that registers event grid system topic for the resource group deployment events
// and a subscription using the provided options
func create(options eventGridRegistrationTaskOptions) tasks.Task {
	action := func(ctx context.Context) error {

		ready := false
		for ready {
			ready = options.IsReady()
			log.Info("Ready: %v", ready)
			time.Sleep(1 * time.Second)
		}

		manager, err := subscriptionmanagement.NewEventGridManager(options.Credential, options.ResourceGroupId)

		if err != nil {
			log.Printf("Error creating event grid manager: %v", err)
			return err
		}
		log.Printf("EventGrid manager created for resource group: %s", options.ResourceGroupId)
		resourceId, err := arm.ParseResourceID(options.ResourceGroupId)
		if err != nil {
			log.Printf("Error parsing resource group id: %v", err)
			return err
		}

		_, err = manager.CreateSystemTopic(ctx)
		if err != nil {
			return err
		}
		log.Printf("System topic created: %s", manager.GetSystemTopicName())

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknownhost"
		}

		subscriptionName := resourceId.ResourceGroupName + "-deployment-events-" + hostname
		result, err := manager.CreateEventSubscription(ctx, subscriptionName, options.EndpointUrl)
		if err != nil {
			return err
		}
		log.Printf("EventGrid subscription created: %s", *result.Name)

		return nil
	}
	return tasks.NewTask("EventGrid Subscription Registration", action)
}
