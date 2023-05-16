package operations

import (
	"context"
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
	log    *log.Entry
}

func (exe *dryRun) Execute(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	retries := uint(invokedOperation.Retries)
	exe.log = log.WithFields(log.Fields{
		"operationId":  invokedOperation.ID,
		"deploymentId": invokedOperation.DeploymentId,
		"retries":      retries,
	})

	err := retry.Do(func() error {
		exe.log.WithField("attempt", invokedOperation.Attempts).Debug("attempting dry run")
		azureDeployment := exe.getAzureDeployment(invokedOperation)

		response, err := exe.dryRun(ctx, azureDeployment)

		if err != nil {
			exe.log.Errorf("Error: %v", err)
			invokedOperation.Status = operation.StatusRunning.String()
			invokedOperation.Attempts = invokedOperation.Attempts + 1
			exe.save(invokedOperation)

			return &RetriableError{Err: err, RetryAfter: 10 * time.Second}
		}

		exe.log.WithField("response", response).Debug("Received dry run response from Azure")

		invokedOperation.Status = operation.StatusSuccess.String()
		invokedOperation.Result = response.DryRunResult
		exe.save(invokedOperation)

		hookMessage := exe.mapToEventHookMessage(invokedOperation, &response.DryRunResult)
		hook.Add(ctx, hookMessage)

		return nil
	},
		retry.Attempts(retries),
	)

	if err != nil {
		exe.log.Errorf("Attempts to retry exceeded. Error: %v", err)
		invokedOperation.Status = operation.StatusFailed.String()
		exe.save(invokedOperation)

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

func (exe *dryRun) mapToEventHookMessage(invokedOperation *data.InvokedOperation, result *deployment.DryRunResult) *events.EventHookMessage {
	data := events.DryRunEventData{
		DeploymentId: int(invokedOperation.DeploymentId),
		OperationId:  invokedOperation.ID,
		Status:       result.Status,
		Attempts:     invokedOperation.Attempts,
		Error:        result.Error,
	}

	message := &events.EventHookMessage{
		Type:   string(events.EventTypeDryRunCompleted),
		Status: operation.StatusSuccess.String(),
		Data:   data,
	}
	message.SetSubject(uint(invokedOperation.DeploymentId), nil)

	return message
}

func (exe *dryRun) getAzureDeployment(operation *data.InvokedOperation) *deployment.AzureDeployment {
	retrieved := &data.Deployment{}
	exe.db.First(&retrieved, operation.DeploymentId)

	deployment := &deployment.AzureDeployment{
		SubscriptionId:    retrieved.SubscriptionId,
		Location:          retrieved.Location,
		ResourceGroupName: retrieved.ResourceGroup,
		DeploymentName:    retrieved.Name,
		Template:          retrieved.Template,
		Params:            operation.Parameters,
	}
	exe.log.Debugf("AzureDeployment: %v", deployment)
	return deployment
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
