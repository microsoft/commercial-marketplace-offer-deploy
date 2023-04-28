package operations

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type dryRunOperation struct {
	db      *gorm.DB
	process DryRunProcessorFunc
	sender  messaging.MessageSender
}

func (h *dryRunOperation) Invoke(operation *data.InvokedOperation) error {
	log.Printf("Inside Invoke for DryRun with an operation of %v", *operation)
	azureDeployment := h.getAzureDeployment(operation)
	response := deployment.DryRun(azureDeployment)
	log.Printf("DryRun response is %v", *response)

	operation.Status = *response.Status
	operation.Result = response.DryRunResult
	operation.UpdatedAt = time.Now().UTC()

	err := h.save(operation)
	if err != nil {
		return err
	}

	eventMsg := h.mapWebHookEventMessage(operation, &response.DryRunResult)
	_ = h.sendEvent(eventMsg)

	return nil
}

func (o *dryRunOperation) sendEvent(eventMessage *events.EventHookMessage) error {
	jsonMsg, err := json.Marshal(eventMessage)
	if err != nil {
		log.Printf("Error marshalling event message: %v", err)
	} else {
		log.Printf("Inside sendEvent for DryRun with a jsonMsg of %v", string(jsonMsg))
	}

	ctx := context.TODO()
	results, err := o.sender.Send(ctx, string(messaging.QueueNameEvents), *eventMessage)
	//results, err := o.sender.Send(ctx, string(messaging.QueueNameEvents), string(jsonMsg))
	if err != nil {
		log.Printf("Error sending event message: %v", err)
		return err
	} else {
		log.Printf("Event message sent successfully")
		log.Printf("Inside sendEvent for DryRun with a results of %v", results)
	}

	if len(results) > 0 {
		for _, result := range results {
			if result.Error != nil {
				log.Printf("Error sending event message: %v", result.Error)
				return result.Error
			}
		}
	}
	return nil
}

func (h *dryRunOperation) mapWebHookEventMessage(operation *data.InvokedOperation, dryRunResult *deployment.DryRunResult) *events.EventHookMessage {
	eventType := "DryRunResult"
	return &events.EventHookMessage{
		Id:        uuid.New(),
		HookId:    [16]byte{},
		EventType: eventType,
		Body:      dryRunResult,
	}
}

func (h *dryRunOperation) getAzureDeployment(operation *data.InvokedOperation) *deployment.AzureDeployment {
	retrieved := &data.Deployment{}
	h.db.First(&retrieved, operation.DeploymentId)

	return &deployment.AzureDeployment{
		SubscriptionId:    retrieved.SubscriptionId,
		Location:          retrieved.Location,
		ResourceGroupName: retrieved.ResourceGroup,
		DeploymentName:    retrieved.Name,
		Template:          retrieved.Template,
		Params:            operation.Parameters,
	}
}

func (h *dryRunOperation) save(operation *data.InvokedOperation) error {
	tx := h.db.Begin()
	tx.Save(&operation)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

//region factory

func NewDryRunProcessor(appConfig *config.AppConfig) DeploymentOperation {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	credential := hosting.GetAzureCredential()
	sender, _ := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
		SubscriptionId:          appConfig.Azure.SubscriptionId,
		Location:                appConfig.Azure.Location,
		ResourceGroupName:       appConfig.Azure.ResourceGroupName,
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	})

	dryRunOperation := &dryRunOperation{
		db:      db,
		process: deployment.DryRun,
		sender:  sender,
	}
	return dryRunOperation
}

//endregion factory
