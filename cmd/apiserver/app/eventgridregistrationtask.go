package app

import (
	"context"
	"net/url"
	"os"
	"regexp"
	"strings"

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
	taskOptions := createOptions(appConfig)
	task := create(taskOptions)

	return task
}

func createOptions(appConfig *config.AppConfig) eventGridRegistrationTaskOptions {
	route, _ := url.Parse("eventgrid")
	return eventGridRegistrationTaskOptions{
		CredentialFunc:  hosting.GetAzureCredentialFunc(),
		ResourceGroupId: appConfig.Azure.GetResourceGroupId(),
		EndpointUrl:     appConfig.GetPublicBaseUrl().ResolveReference(route).String(),
	}
}

type eventGridRegistrationTaskOptions struct {
	CredentialFunc  func() azcore.TokenCredential
	ResourceGroupId string
	EndpointUrl     string
}

// factory for creating task that registers event grid system topic for the resource group deployment events
// and a subscription using the provided options
func create(options eventGridRegistrationTaskOptions) tasks.Task {
	action := func(ctx context.Context) error {
		manager, err := subscriptionmanagement.NewEventGridManager(options.CredentialFunc(), options.ResourceGroupId)

		if err != nil {
			log.Errorf("Error creating event grid manager: %v", err)
			return err
		}
		log.Debugf("EventGrid manager created for resource group: %s", options.ResourceGroupId)
		resourceId, err := arm.ParseResourceID(options.ResourceGroupId)
		if err != nil {
			log.Errorf("Error parsing resource group id: %v", err)
			return err
		}

		_, err = manager.CreateSystemTopic(ctx)
		if err != nil {
			return err
		}
		log.Infof("System topic created: %s", manager.GetSystemTopicName())

		subscriptionName := getSubscriptionName(resourceId.ResourceGroupName)
		result, err := manager.CreateEventSubscription(ctx, subscriptionName, options.EndpointUrl)
		if err != nil {
			log.Errorf("EventGrid subscription [%s] failed with error: %s", subscriptionName, err)
			return err
		}
		log.Debugf("EventGrid subscription created: %s", *result.Name)

		return nil
	}
	return tasks.NewTask("EventGrid Subscription Registration", action)
}

func getSubscriptionName(resourceGroupName string) string {
	reservedPrefixes := []string{
		"Microsoft-", "",
		"EventGrid-", "",
		"System-", "",
	}
	replacer := strings.NewReplacer(reservedPrefixes...)
	prefix := replacer.Replace(resourceGroupName)

	r := regexp.MustCompile("[^a-zA-Z0-9 -]")
	prefix = r.ReplaceAllString(prefix, "")

	suffix := "-events-" + getHostname()

	//https://learn.microsoft.com/en-us/azure/event-grid/subscription-creation-schema
	maxLength := 64
	lengthCheck := len(prefix + suffix)

	if lengthCheck > maxLength {
		prefix = prefix[:maxLength-len(suffix)]
	}
	prefix = strings.TrimSuffix(prefix, "-")
	return prefix + suffix
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknownhost"
	}
	return hostname
}
