package azureevents

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

const (
	// we want to identify the write failure and success
	azureEventTypeResourceWriteFailure = "Microsoft.Resources.ResourceWriteFailure"
	azureEventTypeResourceWriteSuccess = "Microsoft.Resources.ResourceWriteSuccess"
	azureDeploymentResourceType        = "Microsoft.Resources/deployments"
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

// composite that holds references to the event grid event and the azure resource
type ResourceEventSubject struct {
	resourceID *arm.ResourceID
	eventData  *ResourceEventData

	// the event grid event message
	Message *eventgrid.Event

	// the azure azureResource instance
	azureResource *armresources.GenericResource

	// the azure deployment resource instance
	azureDeployment *armresources.DeploymentExtended

	operation *operation.Operation

	// the lookup tags for the resource
	tags deployment.LookupTags
}

func NewResourceEventSubject(eventData *ResourceEventData, eventMessage *eventgrid.Event, azureResource *armresources.GenericResource) (*ResourceEventSubject, error) {
	resourceID, err := arm.ParseResourceID(*azureResource.ID)
	if err != nil {
		// this should never happen
		return nil, err
	}

	return &ResourceEventSubject{
		Message:       eventMessage,
		azureResource: azureResource,
		resourceID:    resourceID,
		eventData:     eventData,
	}, nil
}

// region public methods
func (r *ResourceEventSubject) Resource() *armresources.GenericResource {
	return r.azureResource
}

// the associated operation that this event is related to
//
//	remarks:
//		- for example, this would be the "deploy" operation for a deployment
//		- the resource could still be the Stage but the parent would be the deployment
func (r *ResourceEventSubject) Operation() *operation.Operation {
	return r.operation
}

func (r *ResourceEventSubject) AzureDeployment() *armresources.DeploymentExtended {
	return r.azureDeployment
}

func (r *ResourceEventSubject) ResourceID() *arm.ResourceID {
	return r.resourceID
}

func (r *ResourceEventSubject) EventData() *ResourceEventData {
	return r.eventData
}

func (r *ResourceEventSubject) CorrelationID() string {
	return r.eventData.CorrelationID
}

func (r *ResourceEventSubject) Tags() map[string]*string {
	if r.azureResource == nil {
		return map[string]*string{}
	}
	return r.azureResource.Tags
}

func (r *ResourceEventSubject) LookupTags() deployment.LookupTags {
	return r.tags
}

func (r *ResourceEventSubject) SetLookupTags(tags deployment.LookupTags) {
	r.tags = tags
}

// whether the subject of the event is an azure deployment resource
func (r *ResourceEventSubject) IsResourceTypeDeployment() bool {
	if r.azureResource == nil || r.azureResource.Type == nil || *r.azureResource.Type == "" {
		return false
	}
	return *r.azureResource.Type == azureDeploymentResourceType
}

func (r *ResourceEventSubject) IsAzureDeployment() bool {
	return r.azureDeployment != nil
}

//endregion public methods

func (r *ResourceEventSubject) DeploymentId() (*int, error) {
	if !r.IsAzureDeployment() {
		return nil, errors.New("resource is not a deployment")
	}

	if r.IsParentDeployment() {
		// first look into the tags
		if value, ok := r.azureDeployment.Tags[string(deployment.LookupTagKeyDeploymentId)]; ok {
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

func (r *ResourceEventSubject) IsParentDeployment() bool {
	if !r.IsAzureDeployment() {
		return false
	}
	return strings.HasPrefix(*r.azureResource.Name, deployment.LookupPrefix)
}

func (r *ResourceEventSubject) IsStage() bool {
	if !r.IsAzureDeployment() || r.IsParentDeployment() {
		return false
	}
	return true
}

func (r *ResourceEventSubject) IsFailedStage() bool {
	return r.IsStage() && r.GetStatus() == sdk.StatusFailed.String()
}

// get the modm id of the deployment object, which is the stage id
func (r *ResourceEventSubject) StageId() (*uuid.UUID, error) {
	if r.IsStage() {
		if value, ok := r.azureDeployment.Tags[string(deployment.LookupTagKeyId)]; ok {
			if value != nil && *value != "" {
				id, err := uuid.Parse(*value)
				if err == nil {
					return &id, nil
				}
			}
		}
	}
	return nil, errors.New("resource is not a stage")
}

func (r *ResourceEventSubject) GetStatus() string {
	eventType := *r.Message.EventType
	status := sdk.StatusSuccess.String()

	if strings.Contains(eventType, "Failure") {
		status = sdk.StatusFailed.String()
	}
	return status
}

func (r *ResourceEventSubject) GetType() string {
	if !r.IsAzureDeployment() {
		// for anything we receive that isn't a deployment resource, we default to generic event type
		return string(sdk.EventTypeDeploymentEventReceived)
	}

	//default to deployment event
	prefix := "deployment"
	suffix := "Completed"

	if r.IsStage() {
		prefix = "stage"
	}

	operation := r.Operation()
	if operation != nil {
		if operation.IsRetry() {
			suffix = "Retried"
		}
	}
	return prefix + suffix
}
