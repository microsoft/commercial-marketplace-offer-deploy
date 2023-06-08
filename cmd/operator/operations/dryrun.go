package operations

import (
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type dryRunOperation struct {
	dryRun operation.DryRunFunc
	log    *log.Entry
}

func (exe *dryRunOperation) Do(context operation.ExecutionContext) error {
	azureDeployment, err := exe.getAzureDeployment(context.Operation())
	if err != nil {
		return err
	}

	result, err := exe.dryRun(context.Context(), azureDeployment)
	if err != nil {
		return err
	}
	context.Value(result)

	return err
}

func (exe *dryRunOperation) getAzureDeployment(operation *operation.Operation) (*deployment.AzureDeployment, error) {
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
	exe.log.Debugf("Azure Deployment: %+v", deployment)

	return deployment, nil
}

//region factory

func NewDryRunOperation() operation.OperationFunc {
	dryRunOperation := &dryRunOperation{
		dryRun: deployment.DryRun,
		log:    log.WithField("operation", "dryrun"),
	}
	return dryRunOperation.Do
}

//endregion factory
