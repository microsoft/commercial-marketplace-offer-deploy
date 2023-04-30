package operations

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type dryRun struct {
	db     *gorm.DB
	dryRun DryRunFunc
	sender messaging.MessageSender
}

func (exe *dryRun) Execute(ctx context.Context, operation *data.InvokedOperation) error {
	log.Debug("Inside Invoke for DryRun with an operation of %v", *operation)
	azureDeployment := exe.getAzureDeployment(operation)
	response := deployment.DryRun(azureDeployment)
	log.Debug("DryRun response is %v", *response)

	operation.Status = *response.Status
	operation.Result = response.DryRunResult
	operation.UpdatedAt = time.Now().UTC()

	err := exe.save(operation)
	if err != nil {
		return err
	}

	eventMsg := exe.mapToEventHookMessage(operation, &response.DryRunResult)
	_ = exe.sendEvent(eventMsg)

	return nil
}

func (exe *dryRun) sendEvent(eventMessage *events.EventHookMessage) error {
	ctx := context.TODO()
	results, err := exe.sender.Send(ctx, string(messaging.QueueNameEvents), *eventMessage)
	if err != nil {
		log.Error("Error sending event message: %v", err)
		return err
	} else {
		log.Debug("Event message sent successfully")
		log.Debug("Inside sendEvent for DryRun with a results of %v", results)
	}

	if len(results) > 0 {
		for _, result := range results {
			if result.Error != nil {
				log.Error("Error sending event message: %v", result.Error)
				return result.Error
			}
		}
	}
	return nil
}

func (exe *dryRun) mapToEventHookMessage(operation *data.InvokedOperation, dryRunResult *deployment.DryRunResult) *events.EventHookMessage {
	eventType := "DryRunResult"
	return &events.EventHookMessage{
		Id:        uuid.New(),
		HookId:    [16]byte{},
		EventType: eventType,
		Body:      dryRunResult,
	}
}

func (exe *dryRun) getAzureDeployment(operation *data.InvokedOperation) *deployment.AzureDeployment {
	retrieved := &data.Deployment{}
	exe.db.First(&retrieved, operation.DeploymentId)

	return &deployment.AzureDeployment{
		SubscriptionId:    retrieved.SubscriptionId,
		Location:          retrieved.Location,
		ResourceGroupName: retrieved.ResourceGroup,
		DeploymentName:    retrieved.Name,
		Template:          retrieved.Template,
		Params:            operation.Parameters,
	}
}

func (exe *dryRun) save(operation *data.InvokedOperation) error {
	tx := exe.db.Begin()
	tx.Save(&operation)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

//region factory

func NewDryRunExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	credential := hosting.GetAzureCredential()
	sender, _ := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
		SubscriptionId:          appConfig.Azure.SubscriptionId,
		Location:                appConfig.Azure.Location,
		ResourceGroupName:       appConfig.Azure.ResourceGroupName,
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	})

	dryRunOperation := &dryRun{
		db:     db,
		dryRun: deployment.DryRun,
		sender: sender,
	}
	return dryRunOperation
}

//endregion factory
