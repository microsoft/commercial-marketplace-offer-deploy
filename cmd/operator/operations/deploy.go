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

type deployTask struct {
	retryTask           operation.OperationTask
	operationRepository operation.Repository
	deployStageFactory  *operation.DeployStageOperationFactory
}

// the operation to execute
func (task *deployTask) Run(context operation.ExecutionContext) error {
	operation, err := task.getOperation(context)
	if err != nil {
		return err
	}
	return operation(context)
}

func (task *deployTask) Continue(context operation.ExecutionContext) error {
	operation, err := task.getOperation(context)
	if err != nil {
		return err
	}
	return operation(context)
}

func (task *deployTask) getOperation(context operation.ExecutionContext) (operation.OperationFunc, error) {
	run := task.run
	if context.Operation().IsRetry() { // this is a retry if so
		run = task.retryTask.Run
	}
	return run, nil
}

func (task *deployTask) run(context operation.ExecutionContext) error {
	deployStageOperations, err := task.deployStageFactory.Create(context.Operation())
	if err != nil {
		return err
	}

	azureDeployment := task.mapAzureDeployment(context.Operation(), deployStageOperations)

	// save the built arm template to the operation's attributes so we have a snapshot of what was submitted
	context.Attribute(model.AttributeKeyArmTemplate, azureDeployment.Template)

	deployer, err := task.newDeployer(azureDeployment.SubscriptionId)
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

func (task *deployTask) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
}

func (task *deployTask) mapAzureDeployment(parent *operation.Operation, stages map[uuid.UUID]*operation.Operation) deployment.AzureDeployment {
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

func NewDeployTask(appConfig *config.AppConfig) operation.OperationTask {
	repositoryFactory := operation.NewRepositoryFactory(appConfig)
	repository, err := repositoryFactory()
	if err != nil {
		log.Errorf("Failed to create deploy operation: %v", err)
		return nil
	}

	deployStageFactory := operation.NewDeployStageOperationFactory(repository)

	task := &deployTask{
		retryTask:           NewRetryTask(),
		operationRepository: repository,
		deployStageFactory:  deployStageFactory,
	}
	return task
}
