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
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type dryRun struct {
	db     *gorm.DB
	dryRun DryRunFunc
	sender messaging.MessageSender
}

func (exe *dryRun) Execute(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	log.Debugf("Inside Invoke for DryRun with an operation of %v", *invokedOperation)

	err := retry.Do(func() error {
		status := invokedOperation.Status

		log.Debugf("Inside retry.Do for DryRun with an operation of %v", *invokedOperation)
		azureDeployment := exe.getAzureDeployment(invokedOperation)
		log.Debugf("AzureDeployment is %v", *azureDeployment)

		response, err := exe.dryRun(ctx, azureDeployment)

		if err != nil {
			log.Errorf("Error in DryRun: %v", err)

			invokedOperation.Status = operation.StatusFailed.String()
			invokedOperation.Retries = invokedOperation.Retries + 1
			exe.save(invokedOperation)

			return &RetriableError{Err: err, RetryAfter: 10 * time.Second}
		}

		log.Debugf("DryRun response is %v", *response)
		b, err := json.MarshalIndent(*response, "", "  ")
		if err != nil {
			log.Error(err)
		}
		log.Debugf("unmarshaled DryRunResponse: %v", string(b))
		invokedOperation.Status = *response.Status
		invokedOperation.Result = response.DryRunResult
		invokedOperation.UpdatedAt = time.Now().UTC()

		err = exe.save(invokedOperation)
		if err != nil {
			log.Errorf("Error in DryRun: %v", err)
			return &RetriableError{Err: err, RetryAfter: 10 * time.Second}
		}

		hookMessage := exe.mapToEventHookMessage(status, &response.DryRunResult)
		hook.Add(ctx, hookMessage)

		return nil
	},
		retry.Attempts(uint(invokedOperation.Retries)),
	)

	if err != nil {
		log.Errorf("Error in DryRun - Outside of retry loop: %v", err)
		hookMessage := exe.getFailedEventHookMessage(err, invokedOperation)
		hook.Add(ctx, hookMessage)
	}

	return err
}

func (exe *dryRun) getFailedEventHookMessage(err error, invokedOperation *data.InvokedOperation) *events.EventHookMessage {
	var data interface{}
	if err != nil && len(err.Error()) > 0 {
		data = err.Error()
	} else {
		if invokedOperation != nil && invokedOperation.Result != nil {
			data = invokedOperation.Result
		}
	}
	return &events.EventHookMessage{
		Type:   string(events.EventTypeDryRunCompleted),
		Data:   data,
		Status: operation.StatusFailed.String(),
	}
}

func (exe *dryRun) mapToEventHookMessage(status string, result *deployment.DryRunResult) *events.EventHookMessage {
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
