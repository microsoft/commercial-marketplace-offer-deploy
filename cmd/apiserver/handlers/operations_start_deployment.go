package handlers

import (
	"context"
	"errors"
	"log"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	//"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"

	"gorm.io/gorm"
)

func StartDeployment(deploymentId int, operation api.InvokeDeploymentOperation, db *gorm.DB) (interface{}, error) {
	log.Printf("Inside StartDeployment deploymentId: %d", deploymentId)

	toUpdate := &data.Deployment{}
	db.First(&toUpdate, deploymentId)
	pendingStatus := strings.Replace(events.DeploymentPendingEventType.String(), "deployment.", "", 1)
	caser := cases.Title(language.English)
	toUpdate.Status = caser.String(pendingStatus)
	db.Save(toUpdate)

	templateParams := operation.Parameters
	if templateParams == nil {
		return nil, errors.New("templateParams were not provided")
	}

	ctx := context.TODO()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	// Post message to service bus operator queue
	message := data.InvokedOperation{
		DeploymentId:   uint(deploymentId),
		DeploymentName: *operation.Name,
		Parameters:     templateParams.(map[string]interface{}),
	}

	err = enqueueForPublishing(credential, message, ctx)
	if err != nil {
		return nil, err
	}

	uuid := uuid.New().String()
	status := "OK"
	returnedResult := api.InvokedOperation{
		ID:     &uuid,
		Result: events.DeploymentPendingEventType.String(),
		Status: &status,
	}
	return returnedResult, nil
}

func enqueueForPublishing(credential *azidentity.DefaultAzureCredential, message data.InvokedOperation, ctx context.Context) error {
	// sender, err := getMessageSender(credential)
	// if err != nil {
	// 	return err
	// }
	// results, err := sender.Send(ctx, messaging.OperatorQueue, message)
	// if err != nil {
	// 	return err
	// }

	// if len(results) > 0 {
	// 	return utils.NewAggregateError(getErrorMessages(results))
	// }
	return nil
}
func getMessageSender(credential azcore.TokenCredential) (messaging.MessageSender, error) {
	appConfig := config.GetAppConfig()
	sender, err := messaging.NewServiceBusMessageSender(messaging.MessageSenderOptions{
		SubscriptionId:      appConfig.Azure.SubscriptionId,
		Location:            appConfig.Azure.Location,
		ResourceGroupName:   appConfig.Azure.ResourceGroupName,
		ServiceBusNamespace: appConfig.Azure.ServiceBusNamespace,
	}, credential)
	if err != nil {
		return nil, err
	}

	return sender, nil
}
func getErrorMessages(sendResults []messaging.SendMessageResult) []string {
	errors := []string{}
	for _, result := range sendResults {
		if result.Error != nil {
			errors = append(errors, result.Error.Error())
		}
	}

	return errors
}
