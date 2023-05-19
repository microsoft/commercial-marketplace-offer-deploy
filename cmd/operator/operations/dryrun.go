package operations

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/avast/retry-go"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type dryRun struct {
	db         *gorm.DB
	dryRun     DryRunFunc
	sender     messaging.MessageSender
	log        *log.Entry
	retryDelay time.Duration
}

func (exe *dryRun) Execute(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	retries := uint(invokedOperation.Retries)
	exe.log = log.WithFields(log.Fields{
		"operationId":  invokedOperation.ID,
		"deploymentId": invokedOperation.DeploymentId,
	})

	err := retry.Do(func() error {
		log := exe.log.WithField("attempt", invokedOperation.Attempts+1)
		log.Info("attempting dry run")

		invokedOperation.Running()
		exe.save(invokedOperation)

		azureDeployment := exe.getAzureDeployment(invokedOperation)
		response, err := exe.dryRun(ctx, azureDeployment)

		if err != nil {
			log.Errorf("Error: %v", err)
			invokedOperation.Error(err)
			exe.save(invokedOperation)

			return &RetriableError{Err: err, RetryAfter: exe.retryDelay}
		}
		log.WithField("response", response).Debug("Received dry run response from Azure")

		var result []*sdk.DryRunResult
		if response != nil {
			result = response
			invokedOperation.Status = sdk.StatusSuccess.String()
		} else {
			invokedOperation.Status = string(sdk.StatusError)
			exe.log.Warn("Dry run response is nil")
		}

		invokedOperation.Result = result
		exe.save(invokedOperation)

		hookMessage := exe.mapToEventHookMessage(invokedOperation, result)
		hook.Add(ctx, hookMessage)

		return nil
	},
		retry.Attempts(retries),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			if retriable, ok := err.(*RetriableError); ok {
				return retriable.RetryAfter
			}
			return retry.BackOffDelay(n, err, config)
		}),
	)

	if err != nil {
		exe.log.Errorf("Attempts to retry exceeded. Error: %v", err)
		invokedOperation.Status = sdk.StatusFailed.String()
		exe.save(invokedOperation)

		hookMessage := exe.getFailedEventHookMessage(err, invokedOperation)
		hook.Add(ctx, hookMessage)
	}

	return err
}

func (exe *dryRun) getFailedEventHookMessage(err error, invokedOperation *data.InvokedOperation) *sdk.EventHookMessage {
	data := &sdk.DryRunEventData{
		DeploymentId: int(invokedOperation.DeploymentId),
		OperationId:  invokedOperation.ID,
		Status:       to.Ptr(sdk.StatusFailed.String()),
		Attempts:     invokedOperation.Attempts,
	}

	message := &sdk.EventHookMessage{
		Type:   string(sdk.EventTypeDryRunCompleted),
		Data:   data,
		Status: sdk.StatusFailed.String(),
	}

	if err != nil && len(err.Error()) > 0 {
		message.Error = err.Error()
	}

	return message
}

func getErrorResponse(results []*sdk.DryRunResult) *sdk.DryRunErrorResponse {
	if results == nil {
		return nil
	}

	var errorAdditionalInfo []*sdk.ErrorAdditionalInfo
	for _, result := range results {
		if result.Status != nil && *result.Status == sdk.StatusError.String() {
			if result.Error != nil {
				errorAdditionalInfo = append(errorAdditionalInfo, result.Error.AdditionalInfo...)
			}
		}
	}

	return nil
}

func getStatus(result []*sdk.DryRunResult) *string {
	if result == nil {
		return to.Ptr(sdk.StatusError.String())
	}

	for _, r := range result {
		if r.Status != nil {
			if *r.Status == sdk.StatusError.String() {
				return to.Ptr(sdk.StatusError.String())
			}
		}
	}

	return to.Ptr(sdk.StatusSuccess.String())
}

func (exe *dryRun) mapToEventHookMessage(invokedOperation *data.InvokedOperation, result []*sdk.DryRunResult) *sdk.EventHookMessage {
	resultStatus := to.Ptr(sdk.StatusError.String())
	resultError := &sdk.DryRunErrorResponse{}

	if result != nil {
		resultStatus = getStatus(result)
		resultError = getErrorResponse(result)
	}

	data := sdk.DryRunEventData{
		DeploymentId: int(invokedOperation.DeploymentId),
		OperationId:  invokedOperation.ID,
		Status:       resultStatus,
		Attempts:     invokedOperation.Attempts,
		Error:        resultError,
		StartedAt:    invokedOperation.CreatedAt,
		CompletedAt:  invokedOperation.UpdatedAt,
	}

	message := &sdk.EventHookMessage{
		Type:   string(sdk.EventTypeDryRunCompleted),
		Status: sdk.StatusSuccess.String(),
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
		DeploymentName:    retrieved.GetAzureDeploymentName(),
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
		db:         db,
		dryRun:     deployment.DryRun,
		sender:     sender,
		retryDelay: 5 * time.Second,
	}
	return dryRunOperation
}

//endregion factory
