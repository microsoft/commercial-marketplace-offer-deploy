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
	Watch(id uuid.UUID, options OperationWatcherOptions) (context.Context, error)
}

type operationWatcher struct {
	repository Repository
}

type watchParameters struct {
	OperationWatcherOptions
	ctx    context.Context
	id     uuid.UUID
	ticker *time.Ticker
}

func NewWatcher(repository Repository) OperationWatcher {
	return &operationWatcher{
		repository: repository,
	}
}

func (watcher *operationWatcher) Watch(id uuid.UUID, options OperationWatcherOptions) (context.Context, error) {
	ctx := context.TODO()

	exists := watcher.repository.Any(id)
	if !exists {
		return ctx, fmt.Errorf("failed to start watcher. operation not found for [%s]", id)
	}

	if options.Condition == nil {
		return ctx, fmt.Errorf("failed to start watcher. condition cannot be nil")
	}

	parameters := watchParameters{
		OperationWatcherOptions: options,
		ticker:                  time.NewTicker(options.Frequency),
		ctx:                     ctx,
		id:                      id,
	}
	go watcher.watch(parameters)

	return ctx, nil
}

func (watcher *operationWatcher) watch(params watchParameters) {
	ctx, cancel := context.WithCancel(params.ctx)
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
				cancel()
			}
		case <-ctx.Done(): // if the context is cancelled, externally, then stop
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
