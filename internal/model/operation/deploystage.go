package operation

import (
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

type DeployStageOperationFactory struct {
	repository Repository
}

type DeployStageOperationMap map[uuid.UUID]*Operation

// create a new deploy stage operation for each stage associated to the parent
func (f *DeployStageOperationFactory) Create(parent *Operation) (DeployStageOperationMap, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent operation is nil")
	}

	stageOperations := DeployStageOperationMap{}
	errorMessages := []string{}

	deployment := parent.Deployment()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(deployment.Stages))

	for _, stage := range deployment.Stages {
		go func(stage model.Stage) {
			defer waitGroup.Done()

			configure := func(o *model.InvokedOperation) error {
				o.Name = string(sdk.OperationDeployStage)
				o.ParentID = to.Ptr(parent.ID)
				o.Retries = stage.Retries
				o.Attempts = 0
				o.Status = string(sdk.StatusNone)
				o.DeploymentId = deployment.ID
				o.Parameters = map[string]any{
					string(model.ParameterKeyStageId):            stage.ID,
					string(model.ParameterKeyNestedTemplateName): stage.AzureDeploymentName,
				}

				return nil
			}

			stageOperation, err := f.repository.New(sdk.OperationDeployStage, configure)
			if err != nil {
				log.Errorf("Failed to create deploy stage operation: %v", err)
				errorMessages = append(errorMessages, err.Error())
				return
			}
			err = stageOperation.SaveChanges()
			if err != nil {
				errorMessages = append(errorMessages, err.Error())
				return
			}
			stageOperations[stage.ID] = stageOperation
		}(stage)
	}
	waitGroup.Wait()

	if len(errorMessages) > 0 {
		return nil, utils.NewAggregateError(errorMessages)
	}
	return stageOperations, nil
}

func NewDeployStageOperationFactory(repository Repository) *DeployStageOperationFactory {
	return &DeployStageOperationFactory{
		repository: repository,
	}
}
