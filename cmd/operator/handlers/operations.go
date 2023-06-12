package handlers

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	log "github.com/sirupsen/logrus"
)

type operationMessageHandler struct {
	repository operation.Repository
}

func (h *operationMessageHandler) Handle(message *operation.ExecuteOperationCommand, context messaging.MessageHandlerContext) error {
	h.repository.WithContext(context.Context())
	operation, err := h.repository.First(message.OperationId)
	if err != nil {
		return err
	}
	return operation.Execute()
}

func NewOperationsMessageHandler(appConfig *config.AppConfig) *operationMessageHandler {
	handler := &operationMessageHandler{}
	manager, err := newOperationManager(appConfig)
	if err != nil {
		log.Errorf("Error creating operations message handler: %s", err)
		return nil
	}

	repository, err := operation.NewRepository(manager, operations.NewOperationTaskProvider(appConfig))
	if err != nil {
		log.Errorf("Error creating operations message handler: %s", err)
		return nil
	}

	handler.repository = repository

	return handler
}

func newOperationManager(appConfig *config.AppConfig) (*operation.OperationManager, error) {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()

	scheduler, err := operation.NewSchedulerFromConfig(appConfig)
	if err != nil {
		return nil, err
	}

	service, err := operation.NewManager(db, scheduler, hook.Notify)
	if err != nil {
		return nil, err
	}
	return service, nil
}
