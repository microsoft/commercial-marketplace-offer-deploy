package eventgrid

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

const (
	// we want to identify the write failure and success
	azureEventTypeResourceWriteFailure = "Microsoft.Resources.ResourceWriteFailure"
	azureEventTypeResourceWriteSuccess = "Microsoft.Resources.ResourceWriteSuccess"
)

// ResourceEventData is the data structure for the event grid event
// use only for unmarshalling in order to map to resource
// this is the .Data field of the event grid event when the event type is Microsoft.Resources.ResourceWriteSuccess,
// for example.
type ResourceEventData struct {
	CorrelationID    string `mapstructure:"correlationId"`
	ResourceProvider string `mapstructure:"resourceProvider"`
	ResourceURI      string `mapstructure:"resourceUri"`
	OperationName    string `mapstructure:"operationName"`
	Status           string `mapstructure:"status"`
	SubscriptionID   string `mapstructure:"subscriptionId"`
	TenantID         string `mapstructure:"tenantId"`
}

// list of event grid event resources
type EventGridEventResources []*EventGridEventResource

// combines event grid event message and the resource instance it's related to
type EventGridEventResource struct {
	// the event grid event message
	Message *eventgrid.Event

	// the azure resource instance
	Resource *armresources.GenericResource

	// the azure deployment resource instance
	Deployment *armresources.DeploymentExtended

	// the lookup tags for the resource
	Tags deployment.LookupTags
}

func (r *EventGridEventResource) IsDeployment() bool {
	return r.Deployment != nil
}

func (r *EventGridEventResource) CorrelationID() (*string, error) {
	if !r.IsDeployment() {
		return nil, errors.New("event grid event resource is not a deployment")
	}
	return r.Deployment.Properties.CorrelationID, nil
}

func (r *EventGridEventResource) DeploymentId() (*int, error) {
	if !r.IsDeployment() {
		return nil, errors.New("resource is not a deployment")
	}

	if r.IsParentDeployment() {
		// first look into the tags
		if value, ok := r.Deployment.Tags[string(deployment.LookupTagKeyDeploymentId)]; ok {
			if value == nil || *value == "" {
				return nil, errors.New("deployment id is empty")
			}

			id, err := strconv.Atoi(*value)
			if err == nil {
				return &id, nil
			}
		}
	}
	return nil, errors.New("deployment id not found")
}

func (r *EventGridEventResource) IsParentDeployment() bool {
	if !r.IsDeployment() {
		return false
	}
	return strings.HasPrefix(*r.Resource.Name, deployment.LookupPrefix)
}

func (r *EventGridEventResource) IsStage() bool {
	if !r.IsDeployment() || r.IsParentDeployment() {
		return false
	}
	return true
}

func (r *EventGridEventResource) IsFailedStage() bool {
	return r.IsStage() && r.GetStatus() == sdk.StatusFailed.String()
}

// get the modm id of the deployment object, which is the stage id
func (r *EventGridEventResource) StageId() (uuid.UUID, error) {
	if r.IsStage() {
		if value, ok := r.Deployment.Tags[string(deployment.LookupTagKeyId)]; ok {
			if value != nil && *value != "" {
				id, err := uuid.Parse(*value)
				if err == nil {
					return id, nil
				}
			}
		}
	}
	return uuid.Nil, errors.New("resource is not a stage")
}

func (r *EventGridEventResource) GetStatus() string {
	eventType := *r.Message.EventType

	switch eventType {
	case azureEventTypeResourceWriteFailure:
		return sdk.StatusFailed.String()
	case azureEventTypeResourceWriteSuccess:
		return sdk.StatusSuccess.String()
	default:
		return strcase.ToCamel(eventType)
	}
}
