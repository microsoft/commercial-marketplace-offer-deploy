package operation

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
)

type OperationWatcherOptions struct {
	Condition func(operation model.InvokedOperation) bool // the condition to watch for
	Frequency time.Duration
}

type OperationWatcher interface {
	// watch the operation with the given operation id for completion
	// 	id: the operation id
	//
	Watch(id uuid.UUID, options OperationWatcherOptions) (*OperationWatcherHandle, error)
}

type OperationWatcherHandle struct {
	Context context.Context
	Done    context.CancelFunc
}

type operationWatcher struct {
	repository Repository
}

type watchParameters struct {
	OperationWatcherOptions
	handle *OperationWatcherHandle
	id     uuid.UUID
	ticker *time.Ticker
}

func NewWatcher(repository Repository) OperationWatcher {
	return &operationWatcher{
		repository: repository,
	}
}

func (watcher *operationWatcher) Watch(id uuid.UUID, options OperationWatcherOptions) (*OperationWatcherHandle, error) {
	ctx, cancel := context.WithCancel(context.TODO())

	handle := &OperationWatcherHandle{
		Context: ctx,
		Done:    cancel,
	}

	exists := watcher.repository.Any(id)
	if !exists {
		return handle, fmt.Errorf("failed to start watcher. operation not found for [%s]", id)
	}

	if options.Condition == nil {
		return handle, fmt.Errorf("failed to start watcher. condition cannot be nil")
	}

	parameters := watchParameters{
		OperationWatcherOptions: options,
		ticker:                  time.NewTicker(options.Frequency),
		handle:                  handle,
		id:                      id,
	}
	go watcher.watch(parameters)

	return handle, nil
}

func (watcher *operationWatcher) watch(params watchParameters) {
	ticker := params.ticker

	for {
		select {
		case <-ticker.C:
			operation, err := watcher.getOperation(params.id)
			if err != nil {
				log.Warnf("failed to get operation [%s]. %s", params.id, err)
			}
			evaluation := params.Condition(operation.InvokedOperation)
			if evaluation {
				ticker.Stop()
				params.handle.Done()
			}
		case <-params.handle.Context.Done(): // if the context is cancelled, externally, then stop
			ticker.Stop()
		}
	}
}

func (watcher *operationWatcher) getOperation(id uuid.UUID) (*Operation, error) {
	operation, err := watcher.repository.First(id)
	if err != nil {
		return nil, err
	}

	return operation, nil
}
