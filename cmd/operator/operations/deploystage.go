package operations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/gommon/log"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/stage"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type nameFinderFactory func(context operation.ExecutionContext) (*operation.AzureDeploymentNameFinder, error)

type deployStageOperation struct {
	pollerFactory     *stage.DeployStagePollerFactory
	nameFinderFactory nameFinderFactory
	watcher           operation.OperationWatcher
}

func (op *deployStageOperation) Do(executionContext operation.ExecutionContext) error {
	watcherCtx, err := op.watchParentOperation(executionContext)
	if err != nil {
		return err
	}

	finder, err := op.nameFinderFactory(executionContext)
	if err != nil {
		return err
	}

	azureDeploymentName, err := finder.FindUntilDone(watcherCtx)
	if err != nil {
		return err
	}

	// save the deployment name to the operation so we can fetch it later
	executionContext.Operation().Attribute(model.AttributeKeyAzureDeploymentName, azureDeploymentName)
	executionContext.SaveChanges()

	isFirstAttempt := executionContext.Operation().IsFirstAttempt()
	if isFirstAttempt {
		err := op.wait(executionContext, azureDeploymentName)
		if err != nil {
			return err
		}
	} else { // retry the stage
		retryStage := NewRetryStageOperation()
		err := retryStage(executionContext)
		if err != nil {
			return err
		}
	}
	return nil
}

// watches the parent deploy operation for failure or completed state
// it will trigger a cancellation of the ctx on the execution context if the condition is met
func (op *deployStageOperation) watchParentOperation(context operation.ExecutionContext) (context.Context, error) {
	parentId := context.Operation().ParentID
	if parentId == nil {
		return nil, errors.New("parent operation id is nil")
	}
	options := operation.OperationWatcherOptions{
		Condition: func(operation model.InvokedOperation) bool {
			return operation.Status == sdk.StatusFailed.String() && operation.IsCompleted()
		},
		Frequency: 5 * time.Second,
	}
	ctx, err := op.watcher.Watch(*parentId, options)
	if err != nil {
		return nil, fmt.Errorf("failed to start watcher for parent operation [%s]", *parentId)
	}
	return ctx, nil
}

func (op *deployStageOperation) wait(context operation.ExecutionContext, azureDeploymentName string) error {
	poller, err := op.pollerFactory.Create(context.Operation(), azureDeploymentName, nil)
	if err != nil {
		return err
	}
	response, err := poller.PollUntilDone(context.Context())
	if err != nil {
		return err
	}

	context.Value(response)

	if response.Status == sdk.StatusFailed {
		return operation.NewError(context.Operation())
	}

	return nil
}

func NewDeployStageOperation(appConfig *config.AppConfig) operation.OperationFunc {
	pollerFactory := stage.NewDeployStagePollerFactory()

	repository, err := newOperationRepository(appConfig)
	if err != nil {
		log.Errorf("failed construct deployStage operation: %s", err)
		return nil
	}

	operation := &deployStageOperation{
		watcher:       operation.NewWatcher(repository),
		pollerFactory: pollerFactory,
		nameFinderFactory: func(context operation.ExecutionContext) (*operation.AzureDeploymentNameFinder, error) {
			return operation.NewAzureDeploymentNameFinder(context.Operation())
		},
	}
	return operation.Do
}

func newOperationRepository(appConfig *config.AppConfig) (operation.Repository, error) {
	manager, err := newOperationManager(appConfig)
	if err != nil {
		return nil, err
	}

	repository, err := operation.NewRepository(manager, nil)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func newOperationManager(appConfig *config.AppConfig) (*operation.OperationManager, error) {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
		SubscriptionId:          appConfig.Azure.SubscriptionId,
		Location:                appConfig.Azure.Location,
		ResourceGroupName:       appConfig.Azure.ResourceGroupName,
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	})

	if err != nil {
		return nil, err
	}

	service, err := operation.NewManager(db, sender, hook.Notify)
	if err != nil {
		return nil, err
	}
	return service, nil
}
