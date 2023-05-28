package operations

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type dryRunOperation struct {
	dryRun     operation.DryRunFunc
	log        *log.Entry
	retryDelay time.Duration
}

func (exe *dryRunOperation) Do(context *operation.ExecutionContext) error {
	invokedOperation := context.InvokedOperation()

	retries := uint(invokedOperation.Retries)

	exe.log = log.WithFields(log.Fields{
		"operationId":  invokedOperation.ID,
		"deploymentId": invokedOperation.DeploymentId,
	})

	err := retry.Do(func() error {
		log := exe.log.WithField("attempt", invokedOperation.Attempts+1)

		err := invokedOperation.Running()

		if err != nil {
			return err
		}

		log.Info("attempting dry run")

		azureDeployment := exe.getAzureDeployment(invokedOperation)
		result, err := exe.dryRun(context.Context(), azureDeployment)

		if err != nil {
			log.Errorf("error executing dry run. error: %v", err)
			invokedOperation.Error(err)
			invokedOperation.SaveChanges() //save changes instead of retry since we're doing it inline with retry.Do

			return &operation.RetriableError{Err: err, RetryAfter: exe.retryDelay}
		}

		log.WithField("result", result).Debug("Received dry run result")

		if result != nil {
			invokedOperation.Value(result)

		} else {
			exe.log.Warn("Dry run result was nil")
		}

		invokedOperation.Success()

		return nil
	},
		retry.Attempts(retries),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			if retriable, ok := err.(*operation.RetriableError); ok {
				return retriable.RetryAfter
			}
			return retry.BackOffDelay(n, err, config)
		}),
	)

	if err != nil {
		exe.log.Errorf("Attempts to retry exceeded. Error: %v", err)
	}

	return err
}

func (exe *dryRunOperation) getAzureDeployment(operation *operation.Operation) *deployment.AzureDeployment {
	d := operation.Deployment()

	deployment := &deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		Location:          d.Location,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            operation.Parameters,
	}
	exe.log.Debugf("AzureDeployment: %v", deployment)
	return deployment
}

//region factory

func NewdryRunOperation(appConfig *config.AppConfig) operation.OperationFunc {
	dryRunOperation := &dryRunOperation{
		dryRun:     deployment.DryRun,
		retryDelay: 5 * time.Second,
	}
	return dryRunOperation.Do
}

//endregion factory
