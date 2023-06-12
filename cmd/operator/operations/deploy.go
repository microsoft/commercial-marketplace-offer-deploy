package operations

import (
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/template"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type deployOperation struct {
	retryOperation      operation.OperationFunc
	operationRepository operation.Repository
	deployStageFactory  *operation.DeployStageOperationFactory
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
	deployStageOperations, err := op.deployStageFactory.Create(context.Operation())
	if err != nil {
		return err
	}

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

func (op *deployOperation) mapAzureDeployment(parent *operation.Operation, stages map[uuid.UUID]*operation.Operation) deployment.AzureDeployment {
	d := parent.Deployment()

	template := template.NewDeploymentTemplate(d.Template)

	for _, stage := range d.Stages {
		stageOperation := stages[stage.ID]

		if nestedTemplateName, ok := stageOperation.ParameterValue(model.ParameterKeyNestedTemplateName); ok {
			if value, ok := nestedTemplateName.(string); ok {
				lookupTag := deployment.LookupTag{
					Key:   deployment.LookupTagKeyOperationId,
					Value: to.Ptr(stageOperation.ID.String()),
				}
				template.Tag(value, lookupTag)
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
	repositoryFactory := operation.NewRepositoryFactory(appConfig)
	repository, err := repositoryFactory()
	if err != nil {
		log.Errorf("Failed to create deploy operation: %v", err)
		return nil
	}

	deployStageFactory := operation.NewDeployStageOperationFactory(repository)

	operation := &deployOperation{
		retryOperation:      NewRetryOperation(),
		operationRepository: repository,
		deployStageFactory:  deployStageFactory,
	}
	return operation.Do
}
