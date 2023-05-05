package operations

import (
	"context"
	"encoding/json"
	"time"
	"github.com/avast/retry-go"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
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

	err := retry.Do(func() error {
		log.Debug("Inside retry.Do for DryRun with an operation of %v", *operation)
		azureDeployment := exe.getAzureDeployment(operation)
		log.Debug("AzureDeployment is %v", *azureDeployment)

		response, err := exe.dryRun(ctx, azureDeployment)
		
		if err != nil {
			log.Error("Error in DryRun: %v", err)
			return &RetriableError{Err: err, RetryAfter: 10 * time.Second}
		}

		log.Debug("DryRun response is %v", *response)
		b, err := json.MarshalIndent(*response, "", "  ")
		if err != nil {
			log.Error(err)
		}
		log.Debugf("unmarshaled DryRunResponse: %v", string(b))
		operation.Status = *response.Status
		operation.Result = response.DryRunResult
		operation.UpdatedAt = time.Now().UTC()

		err = exe.save(operation)
		if err != nil {
			log.Error("Error in DryRun: %v", err)
			return &RetriableError{Err: err, RetryAfter: 10 * time.Second}
		}
		
		hookMessage := exe.mapToEventHookMessage(&response.DryRunResult)
		hook.Add(ctx, hookMessage)

		return nil
	},
	retry.Attempts(uint(operation.Retries)),
	)

	return err
}

func (exe *dryRun) mapToEventHookMessage(result *deployment.DryRunResult) *events.EventHookMessage {
	return &events.EventHookMessage{
		Type: string(events.EventTypeDryRunCompleted),
		Data: result,
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
	log.Debug("Inside NewDryRunExecutor")
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
