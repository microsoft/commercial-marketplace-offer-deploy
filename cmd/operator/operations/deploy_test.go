package operations

// import (
// 	"context"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
// 	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
// 	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
// 	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
// 	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
// 	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
// 	"github.com/microsoft/commercial-marketplace-offer-deploy/test/fakes"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/gorm"
// )

// type deployTest struct {
// 	hook                  *fakes.FakeHookQueue
// 	db                    *gorm.DB
// 	createAzureDeployment deployment.CreateDeployment
// 	executionFactory      operation.ExecutorFactory
// 	invokedOperation      *model.InvokedOperation
// 	t                     *testing.T
// }

// func newDeployTest(t *testing.T) *deployTest {
// 	// set hook queue to fake instance
// 	fakeHookQueue := fakes.NewFakeHookQueue(t)
// 	hook.SetInstance(fakeHookQueue)

// 	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
// 	d := model.Deployment{}
// 	db.Create(&d)

// 	createAzureDeployment := func(ctx context.Context, dep deployment.AzureDeployment) (*deployment.AzureDeploymentResult, error) {
// 		t.Log("createAzureDeployment called")
// 		return &deployment.AzureDeploymentResult{}, nil
// 	}

// 	invokedOperation := &model.InvokedOperation{
// 		BaseWithGuidPrimaryKey: model.BaseWithGuidPrimaryKey{ID: uuid.New()},
// 		Name:                   "startDeployment",
// 		DeploymentId:           d.ID,
// 		Attempts:               0,
// 		Retries:                1,
// 	}
// 	db.Create(&invokedOperation)

// 	return &deployTest{
// 		t:                     t,
// 		db:                    db,
// 		hook:                  fakeHookQueue,
// 		executionFactory:      &fakeExecutorFactory{},
// 		createAzureDeployment: createAzureDeployment,
// 		invokedOperation:      invokedOperation,
// 	}
// }

// func Test_Deploy_FirstAttemptSendsEventHookWithOperationId(t *testing.T) {
// 	test := newDeployTest(t)

// 	executor := deployeOperation{
// 		db:                    test.db,
// 		factory:               test.executionFactory,
// 		createAzureDeployment: test.createAzureDeployment,
// 	}

// 	// execute with setup
// 	err := executor.Execute(context.TODO(), test.invokedOperation)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, len(test.hook.Messages()))

// 	message := test.hook.Messages()[0]

// 	data, err := message.DeploymentEventData()
// 	assert.NoError(t, err)

// 	assert.EqualValues(t, sdk.EventTypeDeploymentStarted.String(), message.Type)
// 	assert.EqualValues(t, test.invokedOperation.ID, data.OperationId)
// }

// // region fakes

// type fakeExecutorFactory struct {
// }

// type fakeExecutor struct {
// }

// func (f *fakeExecutor) Execute(context *operation.ExecutionContext) error {
// 	return nil
// }

// func (f *fakeExecutorFactory) Create(operationType sdk.OperationType) (operation.Executor, error) {
// 	return &fakeExecutor{}, nil
// }

// //endregion fakes
