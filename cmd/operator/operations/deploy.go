package operations

import (
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/template"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type deployOperation struct {
	retryOperation      operation.OperationFunc
	operationRepository operation.Repository
}

// the operation to execute
func (op *deployOperation) Do(context operation.ExecutionContext) error {
	operation, err := op.getOperation(context)
	if err != nil {
		return err
	}
	return operation(context)
}

func (op *deployOperation) getOperation(context operation.ExecutionContext) (operation.OperationFunc, error) {
	do := op.do

	if context.Operation().IsRetry() { // this is a retry if so
		do = op.retryOperation
	}
	return do, nil
}

func (op *deployOperation) do(context operation.ExecutionContext) error {

	deployStageOperations := op.createDeployStageOperations(context.Operation())

	azureDeployment := op.mapAzureDeployment(context.Operation(), deployStageOperations)
	deployer, err := op.newDeployer(azureDeployment.SubscriptionId)
	if err != nil {
		return err
	}

	beginResult, err := deployer.Begin(context.Context(), azureDeployment)
	if err != nil {
		return err
	}

	// now schedule the operations for all deployStage operations
	for _, stageOperation := range deployStageOperations {
		go stageOperation.Schedule()
	}

	token := beginResult.ResumeToken

	context.Attribute(model.AttributeKeyResumeToken, token)
	context.Attribute(model.AttributeKeyCorrelationId, *beginResult.CorrelationID)

	result, err := deployer.Wait(context.Context(), &token)
	context.Value(result)

	if err != nil {
		return err
	}

	return nil
}

func (op *deployOperation) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
}

func (op *deployOperation) createDeployStageOperations(parent *operation.Operation) map[uuid.UUID]*operation.Operation {
	stageOperations := make(map[uuid.UUID]*operation.Operation)

	deployment := parent.Deployment()

	for _, stage := range deployment.Stages {
		configure := func(stageOperation *model.InvokedOperation) error {
			stageOperation.Name = string(sdk.OperationDeployStage)
			stageOperation.ParentID = to.Ptr(parent.ID)
			stageOperation.Retries = stage.Retries
			stageOperation.Attempts = 0
			stageOperation.Status = string(sdk.StatusUnknown)
			stageOperation.DeploymentId = deployment.ID
			stageOperation.Parameters = map[string]any{
				string(model.ParameterKeyStageId): stage.ID,
			}

			return nil
		}

		stageOperation, err := op.operationRepository.New(sdk.OperationDeployStage, configure)
		if err != nil {
			log.Errorf("Failed to create deploy stage operation: %v", err)
			continue
		}
		stageOperation.SaveChanges()
		stageOperations[stage.ID] = stageOperation
	}

	return stageOperations
}

func (op *deployOperation) mapAzureDeployment(parent *operation.Operation, stages map[uuid.UUID]*operation.Operation) deployment.AzureDeployment {
	d := parent.Deployment()

	template := template.NewDeploymentTemplate(d.Template)

	for _, stage := range d.Stages {
		stageOperation := stages[stage.ID]

		if azureDeploymentName, ok := stageOperation.ParameterValue(model.ParameterKeyAzureDeploymentName); ok {
			if azureDeploymentNameString, ok := azureDeploymentName.(string); ok {
				lookupTag := deployment.LookupTag{
					Key:   deployment.LookupTagKeyOperationId,
					Value: to.Ptr(stageOperation.ID.String()),
				}
				template.Tag(azureDeploymentNameString, lookupTag)
			}
		}
	}

	builtTemplate := template.Build()

	azureDeployment := deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          builtTemplate,
		Params:            parent.Parameters,
		OperationId:       parent.ID,
		Tags: map[string]*string{
			string(deployment.LookupTagKeyDeploymentId): to.Ptr(strconv.Itoa(int(parent.DeploymentId))),
			string(deployment.LookupTagKeyOperationId):  to.Ptr(parent.ID.String()),
		},
	}

	return azureDeployment
}

func NewDeployOperation(appConfig *config.AppConfig) operation.OperationFunc {
	repository, err := newOperationRepository(appConfig)
	if err != nil {
		log.Errorf("Failed to create deploy operation: %v", err)
		return nil
	}
	operation := &deployOperation{
		retryOperation:      NewRetryOperation(),
		operationRepository: repository,
	}
	return operation.Do
}

func newOperationRepository(appConfig *config.AppConfig) (operation.Repository, error) {
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

	repository, err := operation.NewRepository(service, nil)
	if err != nil {
		return nil, err
	}
	return repository, nil
}
