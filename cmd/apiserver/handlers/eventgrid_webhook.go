package handlers

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/azureevents"
	filter "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/azureevents"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/azureevents/eventhook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/diagnostics/audit"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/notification"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

//region handler

// filter event grid event messages by tags that have events=true
var matchAny deployment.LookupTags = deployment.LookupTags{
	deployment.LookupTagKeyEvents: to.Ptr("true"),
}

type eventGridWebHook struct {
	messageFactory      *eventhook.EventHookMessageFactory
	filter              filter.EventGridEventFilter
	stageQuery          *data.StageQuery
	operationQuery      *data.InvokedOperationQuery
	operationRepository operation.Repository
	auditLog            audit.Log
}

// HTTP handler is the webook endpoint that receives event grid events
// the validation middleware will handle validation requests first before this is reached
func (h *eventGridWebHook) Handle(c echo.Context) error {
	log.Debug("Received event grid webhook")
	ctx := c.Request().Context()

	events := []*eventgrid.Event{}
	err := c.Bind(&events)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.audit(events)

	resources := h.filter.Filter(ctx, matchAny, events)

	if len(resources) == 0 {
		return c.String(http.StatusOK, "OK")
	}

	// handle failed events for retry
	err = h.handleFailedDeployment(ctx, resources)
	if err != nil {
		log.Errorf("Failed to handle failed deployment: %s", err.Error())
	}

	err = h.sendEventHookMessages(ctx, resources)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "OK")
}

func (h *eventGridWebHook) audit(events []*eventgrid.Event) {
	for _, event := range events {
		h.auditLog.Append(event)
	}
}

func (h *eventGridWebHook) handleFailedDeployment(ctx context.Context, resources []*azureevents.ResourceEventSubject) error {
	for _, resource := range resources {
		if !resource.IsAzureDeployment() {
			continue
		}

		if resource.IsFailedStage() {
			log.Debugf("Handling failed stage: %v", resource.AzureDeployment())
			stageId, err := resource.StageId()
			if err != nil {
				log.Errorf("Failed to get stage id: %s", err.Error())
				continue
			}

			log.Debugf("Handling stageId: %v", stageId)

			correlationId := resource.CorrelationID()

			deployment, stage, err := h.stageQuery.Execute(*stageId)
			if err != nil {
				log.Errorf("Failed to get deployment and stage: %s", err.Error())
				continue
			}

			invokedOperation, err := h.operationQuery.First(*stageId, correlationId)
			if err != nil {
				log.Errorf("Failed to get invoked operation: %s", err.Error())
				continue
			}

			exists := invokedOperation != nil
			var operation *operation.Operation

			if exists {
				operation, err = h.operationRepository.First(invokedOperation.ID)
				if err != nil {
					log.Errorf("Failed to dispatch operation [%s]: %s", invokedOperation.ID, err.Error())
				}
			} else {
				operation, err = h.operationRepository.New(sdk.OperationRetryStage, func(i *model.InvokedOperation) error {
					i.Parameters = make(map[string]any)
					i.Parameters[string(model.ParameterKeyStageId)] = stageId
					i.Attribute(model.AttributeKeyCorrelationId, correlationId)
					i.Retries = uint(stage.Retries)
					i.DeploymentId = deployment.ID
					return nil
				})
				if err != nil {
					log.Errorf("Failed to create operation: %s", err.Error())
					continue
				}
			}

			err = operation.Schedule()
			if err != nil {
				log.Errorf("Failed to schedule operation [%s]: %s", operation.ID, err.Error())
			}
		}
	}
	return nil
}

func (h *eventGridWebHook) sendEventHookMessages(ctx context.Context, resources []*azureevents.ResourceEventSubject) error {
	messages := h.messageFactory.Create(ctx, resources)
	log.Debugf("Event Hook messages total: %d", len(messages))

	if len(messages) == 0 {
		return nil
	}

	err := h.add(ctx, messages)
	return err
}

// send these event grid events through our message bus to be processed and published
// to the web hook endpoints that are subscribed to our MODM events
func (h *eventGridWebHook) add(ctx context.Context, messages []*sdk.EventHookMessage) error {
	errors := []string{}
	for _, message := range messages {
		log.Debugf("Adding event hook message: %+v", message)
		_, err := hook.Notify(ctx, message)

		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	return utils.NewAggregateError(errors)
}

//endregion handler

//region factory

func NewEventGridWebHookHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) echo.HandlerFunc {
	log.Printf("Creating event grid webhook handler")

	return func(c echo.Context) error {
		errors := []string{}

		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		messageFactory, err := newWebHookEventMessageFactory(appConfig.Azure.SubscriptionId, db, credential)
		if err != nil {
			errors = append(errors, err.Error())
		}

		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			errors = append(errors, err.Error())
		}

		sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
			SubscriptionId:          appConfig.Azure.SubscriptionId,
			Location:                appConfig.Azure.Location,
			ResourceGroupName:       appConfig.Azure.ResourceGroupName,
			FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
		})

		if err != nil {
			errors = append(errors, err.Error())
		}

		service, err := operation.NewService(db, sender, hook.Notify, notification.NewStageNotifier(db))
		if err != nil {
			errors = append(errors, err.Error())
		}

		repository, err := operation.NewRepository(service, nil)
		if err != nil {
			errors = append(errors, err.Error())
		}

		eventsFilter, err := newEventsFilter(appConfig, credential, repository)
		if err != nil {
			errors = append(errors, err.Error())
		}

		auditLog, err := audit.NewAppendOnlyFileAuditLog(appConfig.GetLoggingOptions("eventgrid").FilePath)
		if err != nil {
			errors = append(errors, err.Error())
		}

		if len(errors) > 0 {
			err = utils.NewAggregateError(errors)
			log.Errorf("Failed to create event grid webhook handler: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		handler := eventGridWebHook{
			messageFactory:      messageFactory,
			filter:              eventsFilter,
			stageQuery:          data.NewStageQuery(db),
			operationQuery:      data.NewInvokedOperationQuery(db),
			operationRepository: repository,
			auditLog:            auditLog,
		}

		return handler.Handle(c)
	}
}

func newWebHookEventMessageFactory(subscriptionId string, db *gorm.DB, credential azcore.TokenCredential) (*eventhook.EventHookMessageFactory, error) {
	client, err := armresources.NewDeploymentsClient(subscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	return eventhook.NewEventHookMessageFactory(client, db), nil
}

func newEventsFilter(appConfig *config.AppConfig, credential azcore.TokenCredential, repo operation.Repository) (filter.EventGridEventFilter, error) {
	// TODO: probably should come from db as configurable at runtime
	includeKeys := []string{
		string(deployment.LookupTagKeyEvents),
		string(deployment.LookupTagKeyId),
		string(deployment.LookupTagKeyName),
		string(deployment.LookupTagKeyStageId),
	}
	resourceClient, err := filter.NewAzureResourceClient(appConfig.Azure.SubscriptionId, credential)
	if err != nil {
		return nil, err
	}

	provider := filter.NewResourceEventSubjectFactory(resourceClient, repo)
	filter := filter.NewTagsFilter(includeKeys, provider)
	return filter, nil
}

//endregion factory
