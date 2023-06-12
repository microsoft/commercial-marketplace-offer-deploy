package operations

import (
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type dryRunTask struct {
	dryRun operation.DryRunFunc
	log    *log.Entry
}

func (task *dryRunTask) Run(context operation.ExecutionContext) error {
	azureDeployment, err := task.getAzureDeployment(context.Operation())
	if err != nil {
		return err
	}

	result, err := task.dryRun(context.Context(), azureDeployment)
	if err != nil {
		return err
	}
	context.Value(result)

	return err
}

func (task *dryRunTask) Continue(context operation.ExecutionContext) error {
	return nil
}

func (task *dryRunTask) getAzureDeployment(operation *operation.Operation) (*deployment.AzureDeployment, error) {
	d := operation.Deployment()

	if d == nil {
		return nil, fmt.Errorf("deployment [%d] not found for operation: %s", operation.DeploymentId, operation.Name)
	}

	deployment := &deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		Location:          d.Location,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            operation.Parameters,
	}
	task.log.Debugf("Azure Deployment: %+v", deployment)

	return deployment, nil
}

//region factory

func NewDryRunTask() operation.OperationTask {
	return &dryRunTask{
		dryRun: deployment.DryRun,
		log:    log.WithField("operation", "dryrun"),
	}
}

//endregion factory
