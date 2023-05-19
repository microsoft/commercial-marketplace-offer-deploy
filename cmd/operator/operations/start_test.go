package operations

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/fakes"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type startDeploymentTest struct {
	hook                  *fakes.FakeHookQueue
	db                    *gorm.DB
	createAzureDeployment deployment.CreateDeployment
	executionFactory      ExecutorFactory
	invokedOperation      *data.InvokedOperation
	t                     *testing.T
}

func newStartDeploymentTest(t *testing.T) *startDeploymentTest {
	// set hook queue to fake instance
	fakeHookQueue := fakes.NewFakeHookQueue(t)
	hook.SetInstance(fakeHookQueue)

	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	d := data.Deployment{}
	db.Create(&d)

	createAzureDeployment := func(ctx context.Context, dep deployment.AzureDeployment) (*deployment.AzureDeploymentResult, error) {
		t.Log("createAzureDeployment called")
		return &deployment.AzureDeploymentResult{}, nil
	}

	invokedOperation := &data.InvokedOperation{
		BaseWithGuidPrimaryKey: data.BaseWithGuidPrimaryKey{ID: uuid.New()},
		Name:                   "startDeployment",
		DeploymentId:           d.ID,
		Attempts:               0,
		Retries:                1,
	}
	db.Create(&invokedOperation)

	return &startDeploymentTest{
		t:                     t,
		db:                    db,
		hook:                  fakeHookQueue,
		executionFactory:      &fakeExecutorFactory{},
		createAzureDeployment: createAzureDeployment,
		invokedOperation:      invokedOperation,
	}
}

func Test_StartDeployment_FirstAttemptSendsEventHookWithOperationId(t *testing.T) {
	test := newStartDeploymentTest(t)

	executor := startDeployment{
		db:                    test.db,
		factory:               test.executionFactory,
		createAzureDeployment: test.createAzureDeployment,
	}

	// execute with setup
	err := executor.Execute(context.TODO(), test.invokedOperation)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(test.hook.Messages()))

	message := test.hook.Messages()[0]

	data, err := message.DeploymentEventData()
	assert.NoError(t, err)

	assert.EqualValues(t, sdk.EventTypeDeploymentStarted.String(), message.Type)
	assert.EqualValues(t, test.invokedOperation.ID, data.OperationId)
}

// region fakes

type fakeExecutorFactory struct {
}

type fakeExecutor struct {
}

func (f *fakeExecutor) Execute(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	return nil
}

func (f *fakeExecutorFactory) Create(operationType sdk.OperationType) (Executor, error) {
	return &fakeExecutor{}, nil
}

//endregion fakes
