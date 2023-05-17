package operations

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/fakes"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testStartDeploymentExecutor struct {
	hook                  *fakes.FakeHookQueue
	db                    *gorm.DB
	createAzureDeployment deployment.CreateDeployment
	executionFactory      ExecutorFactory
	invokedOperation      *data.InvokedOperation
	t                     *testing.T
}

func newTestStartDeploymentExecutor(t *testing.T) *testStartDeploymentExecutor {
	// set hook queue to fake instance
	fakeHookQueue := fakes.NewFakeHookQueue(t)
	hook.SetInstance(fakeHookQueue)

	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	d := data.Deployment{}
	db.Create(&d)

	createAzureDeployment := func(ctx context.Context, dep deployment.AzureDeployment) (*deployment.AzureDeploymentResult, error) {
		t.Log("createAzureDeployment called")
		return nil, nil
	}

	invokedOperation := &data.InvokedOperation{
		BaseWithGuidPrimaryKey: data.BaseWithGuidPrimaryKey{ID: uuid.New()},
		Name:                   "startDeployment",
		DeploymentId:           d.ID,
		Attempts:               0,
		Retries:                1,
	}
	db.Create(&invokedOperation)

	return &testStartDeploymentExecutor{
		t:                     t,
		db:                    db,
		hook:                  fakeHookQueue,
		executionFactory:      &fakeExecutorFactory{},
		createAzureDeployment: createAzureDeployment,
		invokedOperation:      invokedOperation,
	}
}

func Test_StartDeployment_FirstAttemptSendsEventHookWithOperationId(t *testing.T) {
	test := newTestStartDeploymentExecutor(t)

	executor := startDeployment{
		db:                    test.db,
		factory:               test.executionFactory,
		createAzureDeployment: test.createAzureDeployment,
	}
	// execute with setup
	err := executor.Execute(context.TODO(), test.invokedOperation)
	assert.NoError(t, err)
	message := test.hook.Messages()[0]

	assert.EqualValues(t, test.invokedOperation.ID, message.Data.(*events.DeploymentEventData).OperationId)
}

// region fakes

type fakeExecutorFactory struct {
}

type fakeExecutor struct {
}

func (f *fakeExecutor) Execute(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	return nil
}

func (f *fakeExecutorFactory) Create(operationType operation.OperationType) (Executor, error) {
	return &fakeExecutor{}, nil
}

//endregion fakes
