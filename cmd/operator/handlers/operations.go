package handlers

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	ops "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

type operationMessageHandler struct {
	database data.Database
	factory  ops.DeploymentOperationFactory
}

func (h *operationMessageHandler) Handle(message *messaging.ExecuteInvokedOperation, context messaging.MessageHandlerContext) error {
	db := h.database.Instance()
	var invokedOperation *data.InvokedOperation
	db.First(&invokedOperation, uuid.MustParse(message.OperationId))

	log.Debug("operation id: %s", message.OperationId)
	log.Debug("Invoked Operation from DB: %v", invokedOperation)

	operationType, err := operations.Type(invokedOperation.Name)
	if err != nil {
		log.Error("Error getting operation type: ", err)
	}

	operation, err := h.factory.Create(operationType)
	if err != nil {
		log.Error("Error creating operation: ", err)
		return err
	}
	operationJson, err := json.Marshal(operation)
	if err != nil {
		log.Error("Error marshalling operation: ", err)
	} else {
		log.Debug("Pulled the operation - Operation: %v", string(operationJson))
	}

	return operation.Execute(context.Context(), invokedOperation)
}

func NewOperationsMessageHandler(appConfig *config.AppConfig) *operationMessageHandler {
	return &operationMessageHandler{
		database: data.NewDatabase(appConfig.GetDatabaseOptions()),
		factory:  ops.NewDeploymentOperationFactory(appConfig),
	}
}
