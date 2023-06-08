package operations

import (
	"errors"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/notification"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

type deployeOperation struct {
	retryOperation operation.OperationFunc
	stageNotifier  notification.StageNotifier
}

// the operation to execute
func (op *deployeOperation) Do(context operation.ExecutionContext) error {
	operation, err := op.getOperation(context)
	if err != nil {
		return err
	}
	return operation(context)
}

func (op *deployeOperation) getOperation(context operation.ExecutionContext) (operation.OperationFunc, error) {
	do := op.do

	if context.Operation().IsRetry() { // this is a retry if so
		do = op.retryOperation
	}
	return do, nil
}

func (op *deployeOperation) do(context operation.ExecutionContext) error {
	azureDeployment := op.mapAzureDeployment(context.Operation())
	deployer, err := op.newDeployer(azureDeployment.SubscriptionId)
	if err != nil {
		return err
	}

	beginResult, err := deployer.Begin(context.Context(), azureDeployment)
	if err != nil {
		return err
	}

	token := beginResult.ResumeToken

	context.Attribute(model.AttributeKeyResumeToken, token)
	context.Attribute(model.AttributeKeyCorrelationId, *beginResult.CorrelationID)

	op.notifyForStages(context)

	result, err := deployer.Wait(context.Context(), &token)
	context.Value(result)

	if err != nil {
		return err
	}

	return nil
}

func (service *deployeOperation) notifyForStages(context operation.ExecutionContext) error {
	operation := context.Operation()

	log.WithFields(log.Fields{
		"attempts": operation.Attempts,
		"status":   operation.Status,
	}).Info("notifying for stages")

	op := operation.InvokedOperation

	log.WithFields(log.Fields{
		"attributesLength": len(op.Attributes),
	}).Info("operation attributes information")

	for _, attr := range op.Attributes {
		//log attr key and value
		log.WithFields(log.Fields{
			"key":   attr.Key,
			"value": attr.Value,
		}).Info("operation attribute")
	}

	// only handle scheduled deployment operations, where we want to notify of stages getting scheduled
	if !op.IsRunning() && !op.IsFirstAttempt() {
		return errors.New("not a running deployment operation or not first attempt")
	}

	correlationId, err := op.CorrelationId()
	if err != nil {
		return err
	}

	deployment := operation.Deployment()
	if deployment == nil {
		return errors.New("deployment not found")
	}

	notification := &model.StageNotification{
		OperationId:       op.ID,
		CorrelationId:     *correlationId,
		ResourceGroupName: deployment.ResourceGroup,
		Entries:           []model.StageNotificationEntry{},
	}

	for _, stage := range deployment.Stages {
		notification.Entries = append(notification.Entries, model.StageNotificationEntry{
			StageId: stage.ID,
			Message: sdk.EventHookMessage{
				Id:     uuid.New(),
				Type:   string(sdk.EventTypeStageStarted),
				Status: sdk.StatusRunning.String(),
				Data: sdk.DeploymentEventData{
					EventData: sdk.EventData{
						DeploymentId: int(deployment.ID),
						OperationId:  op.ID,
						Attempts:     1,
						StartedAt:    time.Now().UTC(),
					},
					StageId:       &stage.ID,
					CorrelationId: correlationId,
				},
			},
		})
	}

	err = service.stageNotifier.Notify(context.Context(), notification)
	return err
}

func (op *deployeOperation) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
}

func (op *deployeOperation) mapAzureDeployment(invokedOperation *operation.Operation) deployment.AzureDeployment {
	d := invokedOperation.Deployment()

	return deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            invokedOperation.Parameters,
		OperationId:       invokedOperation.ID,
		Tags: map[string]*string{
			string(deployment.LookupTagKeyDeploymentId): to.Ptr(strconv.Itoa(int(invokedOperation.DeploymentId))),
			string(deployment.LookupTagKeyOperationId):  to.Ptr(invokedOperation.ID.String()),
		},
	}
}

func NewDeploymentOperation(appConfig *config.AppConfig) operation.OperationFunc {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()

	operation := &deployeOperation{
		retryOperation: NewRetryOperation(),
		stageNotifier:  notification.NewStageNotifier(db),
	}
	return operation.do
}
