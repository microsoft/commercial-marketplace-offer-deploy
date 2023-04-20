package app

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/subscriptionmanagement"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
)

// constructor for creating task that registers event grid system topic for the resource group deployment events
func newEventGridRegistrationTask(appConfig *config.AppConfig) tasks.Task {
	taskOptions := eventGridRegistrationTaskOptions{
		Credential:      hosting.GetAzureCredential(),
		ResourceGroupId: appConfig.Azure.GetResourceGroupId(),
		EndpointUrl:     appConfig.Http.GetBaseUrl() + "/eventgrid",
	}
	task := create(taskOptions)

	return task
}

type eventGridRegistrationTaskOptions struct {
	Credential      azcore.TokenCredential
	ResourceGroupId string
	EndpointUrl     string
}

// factory for creating task that registers event grid system topic for the resource group deployment events
// and a subscription using the provided options
func create(options eventGridRegistrationTaskOptions) tasks.Task {
	action := func(ctx context.Context) error {
		log.Print("Registering event grid system topic for resource group deployment events...")

		manager, err := subscriptionmanagement.NewEventGridManager(options.Credential, options.ResourceGroupId)

		if err != nil {
			return err
		}
		resourceId, err := arm.ParseResourceID(options.ResourceGroupId)
		if err != nil {
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
