package operation

type OperationTask interface {
	Run(context ExecutionContext) error
	Continue(context ExecutionContext) error
}

type operationTask struct {
	runFunc      OperationFunc
	continueFunc OperationFunc
}

func (task *operationTask) Run(context ExecutionContext) error {
	return task.runFunc(context)
}

func (task *operationTask) Continue(context ExecutionContext) error {
	return task.continueFunc(context)
}

func NewOperationTask(options OperationTaskOptions) OperationTask {
	return &operationTask{
		runFunc:      options.Run,
		continueFunc: options.Continue,
	}
}

// Represents the task that the operation will perform
type OperationTaskOptions struct {
	Run      OperationFunc //default run function
	Continue OperationFunc //if the run is intereupted, this function will be invoked
}
