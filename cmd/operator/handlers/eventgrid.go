package handlers

import (
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/labstack/echo"
	internal "github.com/microsoft/commercial-marketplace-offer-deploy/internal/azure/eventgrid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/events"
	"gorm.io/gorm"
)

// defines handler for event grid webhook
type EventGridWebHook struct {
	db      *gorm.DB
	mapper  internal.EventGridEventMapper
	filter  events.ResourceFilter
	Handler echo.HandlerFunc
}

// this API handler is the webook endpoint that receives event grid events
// the validation middleware will handle validation requests first before this is reached
func (h *EventGridWebHook) EventGridWebHookHandler(c echo.Context) error {
	messages := []*eventgrid.Event{}
	err := c.Bind(&messages)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	resources := h.mapper.Map(ctx, messages)
	resources = h.filter.Filter(resources)

	if len(resources) == 0 {
		return c.String(http.StatusOK, "No resources to process")
	}

	for _, resource := range resources {
		id := *resource.ID
		tag := *resource.Tags["modm.events"]
		log.Printf("Resource [%s] capture events '%s'", id, tag)
	}

	return c.String(http.StatusOK, "OK")
}

// handler constructor for event grid webhook handler
func NewEventGridWebHook(databaseOptions *data.DatabaseOptions) *EventGridWebHook {
	d := data.NewDatabase(databaseOptions)
	tags := map[string]string{"modm.events": "true"}

	webhook := &EventGridWebHook{
		db:     d.Instance(),
		mapper: getMapper(),
		filter: events.NewResourceTagFilter(tags),
	}
	webhook.Handler = webhook.EventGridWebHookHandler

	return webhook
}

func getMapper() internal.EventGridEventMapper {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil
	}
	return internal.NewEventGridEventMapper(credential)
}
