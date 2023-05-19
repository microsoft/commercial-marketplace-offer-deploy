package operations

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/fakes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

//region test setup

type dryRunExecutorTest struct {
	db               *gorm.DB
	dryRun           DryRunFunc
	sender           messaging.MessageSender
	hookQueue        *fakes.FakeHookQueue
	invokedOperation *data.InvokedOperation
	ctx              context.Context
}

type dryRunExecutorTestOptions struct {
	causeDryRunError            bool
	causeDryRunResultToBeNil    bool
	causeDryRunStatusToBeFailed bool
}

func (t *dryRunExecutorTest) getSavedState() data.InvokedOperation {
	invokedOperation := &data.InvokedOperation{}
	t.db.First(invokedOperation, t.invokedOperation.ID)
	return *invokedOperation
}

func newDryExecutorTest(t *testing.T, options *dryRunExecutorTestOptions) *dryRunExecutorTest {
	log.SetLevel(log.ErrorLevel)

	hookQueue := fakes.NewFakeHookQueue(t)
	hook.SetInstance(hookQueue)

	dryRunFunc := func(ctx context.Context, ad *deployment.AzureDeployment) (*sdk.DryRunResponse, error) {
		t.Log("dryRunFunc called")
		if options.causeDryRunError {
			return nil, errors.New("dryRunFunc error")
		}

		if options.causeDryRunResultToBeNil {
			return nil, nil
		}

		if options.causeDryRunStatusToBeFailed {
			return &sdk.DryRunResponse{
				DryRunResult: sdk.DryRunResult{
					Status: to.Ptr(sdk.StatusFailed.String()),
				},
			}, nil
		}

		return &sdk.DryRunResponse{
			DryRunResult: sdk.DryRunResult{},
		}, nil
	}

	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()

	deployment := &data.Deployment{}
	db.Save(deployment)

	invokedOperation := &data.InvokedOperation{
		BaseWithGuidPrimaryKey: data.BaseWithGuidPrimaryKey{
			ID: uuid.New(),
		},
		DeploymentId: deployment.ID,
		Name:         string(sdk.OperationDryRun),
		Retries:      3,
		Attempts:     0,
		Parameters:   map[string]interface{}{},
		Result:       nil,
	}

	return &dryRunExecutorTest{
		db:               db,
		dryRun:           dryRunFunc,
		sender:           fakes.NewFakeMessageSender(t),
		hookQueue:        hookQueue,
		ctx:              context.Background(),
		invokedOperation: invokedOperation,
	}
}

//endregion test setup

//region Execute

func Test_DryRun_Execute_failure_hook_message_data_is_DryRunEventData(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{
		causeDryRunError:            true,
		causeDryRunResultToBeNil:    false,
		causeDryRunStatusToBeFailed: true,
	})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
		log:        &log.Entry{},
	}

	executor.Execute(test.ctx, test.invokedOperation)

	assert.Equal(t, 1, len(test.hookQueue.Messages()))

	msg := test.hookQueue.Messages()[0]
	data, err := msg.DryRunEventData()

	t.Logf("data: %v", data)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, sdk.StatusFailed.String(), *data.Status)
}

func Test_DryRun_Execute_DryRunError_Returns_Error(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{causeDryRunError: true})

	executor := &dryRun{
		db:     test.db,
		dryRun: test.dryRun,
		sender: test.sender,
	}

	err := executor.Execute(test.ctx, test.invokedOperation)
	assert.Error(t, err, "dryRunFunc error")
}

func Test_DryRun_Execute_DryRunError_Attempts_Equal_Retries_With_Failure(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{causeDryRunError: true})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
	}

	executor.Execute(test.ctx, test.invokedOperation)
	assert.Equal(t, test.invokedOperation.Retries, test.invokedOperation.Attempts)
}

func Test_DryRun_Execute_DryRunError_InvokedOperation_Attempts_Persisted(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{causeDryRunError: true})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
	}

	executor.Execute(test.ctx, test.invokedOperation)

	invokedOperation := test.getSavedState()

	assert.Equal(t, test.invokedOperation.ID, invokedOperation.ID)
	assert.Equal(t, test.invokedOperation.Retries, invokedOperation.Attempts)
}

func Test_DryRun_Execute_DryRunError_Status_Is_Failed(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{causeDryRunError: true})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
	}

	executor.Execute(test.ctx, test.invokedOperation)

	invokedOperation := test.getSavedState()

	assert.Equal(t, sdk.StatusFailed.String(), test.invokedOperation.Status)
	assert.Equal(t, sdk.StatusFailed.String(), invokedOperation.Status)
}

func Test_DryRun_Execute_NoError_With_Nil_Result_Status_Is_Error(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{
		causeDryRunError:         false,
		causeDryRunResultToBeNil: true,
	})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
	}

	executor.Execute(test.ctx, test.invokedOperation)
	assert.Equal(t, sdk.StatusError.String(), test.invokedOperation.Status)
}

func Test_DryRun_Execute_NoError_With_Result_Status_Is_Success(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{
		causeDryRunError:         false,
		causeDryRunResultToBeNil: false,
	})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
	}

	executor.Execute(test.ctx, test.invokedOperation)
	assert.Equal(t, sdk.StatusSuccess.String(), test.invokedOperation.Status)
}

//endregion Execute

func Test_DryRun_getAzureDeployment_name_is_correctly_set(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{
		causeDryRunError:         false,
		causeDryRunResultToBeNil: false,
	})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
		log:        log.WithField("test", "Test_DryRun_getAzureDeployment_name_is_correctly_set"),
	}

	// set the name using the invoked operation's deployment id
	deployment := &data.Deployment{}
	test.db.First(deployment, test.invokedOperation.DeploymentId)

	deployment.Name = "test-deployment/with some slashes-*&^%$#@!_+=.:'\""
	test.db.Save(deployment)

	result := executor.getAzureDeployment(test.invokedOperation)
	assert.NotNil(t, result)
	assert.Equal(t, "modm.1.test-deploymentwith-some-slashes", result.DeploymentName)
}

func Test_DryRun_Execute_failure_captures_errors(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{
		causeDryRunError:            true,
		causeDryRunResultToBeNil:    false,
		causeDryRunStatusToBeFailed: false,
	})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
		log:        &log.Entry{},
	}

	executor.Execute(test.ctx, test.invokedOperation)

	assert.Equal(t, 3, len(test.invokedOperation.Errors))
}

func Test_DryRun_Execute_eventhook_message_attempts_nonzero_index(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{
		causeDryRunError:            true,
		causeDryRunResultToBeNil:    false,
		causeDryRunStatusToBeFailed: false,
	})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
		log:        &log.Entry{},
	}

	executor.Execute(test.ctx, test.invokedOperation)

	messages := test.hookQueue.Messages()

	for index, message := range messages {
		data, _ := message.DryRunEventData()
		assert.NotEqual(t, index+1, data.Attempts)
	}
}

func Test_DryRun_Execute_eventhook_message_times_match_invokedoperation(t *testing.T) {
	test := newDryExecutorTest(t, &dryRunExecutorTestOptions{
		causeDryRunError:            false,
		causeDryRunResultToBeNil:    false,
		causeDryRunStatusToBeFailed: false,
	})

	executor := &dryRun{
		db:         test.db,
		dryRun:     test.dryRun,
		sender:     test.sender,
		retryDelay: 0 * time.Second,
		log:        &log.Entry{},
	}

	executor.Execute(test.ctx, test.invokedOperation)

	assert.Equal(t, 1, len(test.hookQueue.Messages()))

	message := test.hookQueue.Messages()[0]
	t.Logf("message: %+v", message)

	data, err := message.DryRunEventData()
	assert.NoError(t, err)

	assert.Equal(t, test.invokedOperation.CreatedAt, data.StartedAt)
	assert.Equal(t, test.invokedOperation.UpdatedAt, data.CompletedAt)
}
