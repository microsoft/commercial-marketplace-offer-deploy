package handlers

import (
	"context"

	//"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"

	"gorm.io/gorm"
)

func StartDeployment(deploymentId int, operation api.InvokeDeploymentOperation, db *gorm.DB) (*api.InvokedOperation, error) {
	deployment := &data.Deployment{}
	db.First(&deployment, deploymentId)
	deployment.Status = events.DeploymentPendingEventType.String()
	db.Save(&deployment)

	invokedOperation := data.InvokedOperation{
		DeploymentId: uint(deploymentId),
		Name:         *operation.Name,
		Params:       operation.Parameters.(map[string]interface{}),
	}
	invokedOperation.ID = uuid.New()
	db.Save(&invokedOperation)

	err := enqueueForPublishing(&invokedOperation)
	if err != nil {
		return nil, err
	}

	returnedResult := &api.InvokedOperation{
		ID:     to.Ptr(invokedOperation.ID.String()),
		Result: "Accepted",
		Status: &deployment.Status,
	}
	return returnedResult, nil
}

func enqueueForPublishing(invokedOperation *data.InvokedOperation) error {
	ctx := context.TODO()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	sender, err := getMessageSender(credential)
	if err != nil {
		return err
	}
	message := messaging.InvokedOperationMessage{
		OperationId: invokedOperation.ID.String(),
	}
	results, err := sender.Send(ctx, string(messaging.QueueNameOperations), message)
	if err != nil {
		return err
	}
	if len(results) == 1 && !results[0].Success {
		return results[0].Error
	}
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
