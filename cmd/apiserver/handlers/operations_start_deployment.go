package handlers

import (
	"context"
	"errors"
	"log"
	"time"

	//"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"

	"gorm.io/gorm"
)

// test struct
type Message struct {
	DeploymentId int
	OperationId  string
}

func StartDeployment(deploymentId int, operation api.InvokeDeploymentOperation, db *gorm.DB) (interface{}, error) {
	log.Printf("Inside StartDeployment deploymentId: %d", deploymentId)

	toUpdate := &data.Deployment{}
	db.First(&toUpdate, deploymentId)
	db.Model(&toUpdate).Update("status", "Pending") // TODO: update with deployment.Pending

	templateParams := operation.Parameters
	if templateParams == nil {
		return nil, errors.New("templateParams were not provided")
	}

	// ctx := context.TODO()

	// credential, err := azidentity.NewDefaultAzureCredential(nil)
	// if err != nil {
	// 	return nil, err
	// }

	// // Post message to service bus operator queue
	// message := Message{
	// 	DeploymentId: deploymentId,
	// 	OperationId:  *operation.Name,
	// }

	// err = enqueueForPublishing(credential, message, ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// formulate the response

	timestamp := time.Now().UTC()
	status := "OK"

	returnedResult := api.InvokedOperation{
		ID:         operation.Name,
		InvokedOn:  &timestamp,
		Name:       operation.Name,
		Parameters: templateParams.(map[string]interface{}),
		Result:     "0",
		Status:     &status,
	}
	return returnedResult, nil
}

func enqueueForPublishing(credential *azidentity.DefaultAzureCredential, message Message, ctx context.Context) error {
	sender, err := getMessageSender(credential)
	if err != nil {
		return err
	}
	results, err := sender.Send(ctx, messaging.OperatorQueue, message)
	if err != nil {
		return err
	}

	if len(results) > 0 {
		return utils.NewAggregateError(getErrorMessages(results))
	}
	return nil
}
func getMessageSender(credential azcore.TokenCredential) (messaging.MessageSender, error) {
	var appSettings config.AppSettings = config.GetAppSettings()
	sender, err := messaging.NewServiceBusMessageSender(messaging.MessageSenderOptions{
		SubscriptionId:      appSettings.Azure.SubscriptionId,
		Location:            appSettings.Azure.Location,
		ResourceGroupName:   appSettings.Azure.ResourceGroupName,
		ServiceBusNamespace: appSettings.Azure.ServiceBusNamespace,
	}, credential)
	if err != nil {
		return nil, err
	}

	return sender, nil
}
func getErrorMessages(sendResults []messaging.SendMessageResult) *[]string {
	errors := []string{}
	for _, result := range sendResults {
		if result.Error != nil {
			errors = append(errors, result.Error.Error())
		}
	}

	return &errors
}
