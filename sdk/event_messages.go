package sdk

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk/internal"
)

const (
	EventDataTypePrefixDeployment = "deployment"
	EventDataTypePrefixStage      = "stage"
	EventDataTypePrefixDryRun     = "dryRun"
)

// subscription model for MODM webhook events
type EventHookMessage struct {
	// the ID of the message
	Id uuid.UUID `json:"id,omitempty"`

	// the ID of the hook
	HookId uuid.UUID `json:"hookId,omitempty"`

	// the type of the event, .e.g. "dryRunCompleted" found in pkg/events
	Type string `json:"type,omitempty"`

	// the status of the event, e.g. "success"
	Status string `json:"status,omitempty"`

	Error string `json:"error,omitempty"`

	// subject is in format like /deployments/{deploymentId}/stages/{stageId}/operations/{operationName}
	// /deployments/{deploymentId}/operations/{operationName}
	Subject string `json:"subject,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// Event data for a message
type EventData struct {
	DeploymentId int        `json:"deploymentId" mapstructure:"deploymentId"`
	OperationId  uuid.UUID  `json:"operationId" mapstructure:"operationId"`
	Attempts     int        `json:"attempts" mapstructure:"attempts"`
	ScheduledAt  *time.Time `json:"scheduledAt,omitempty" mapstructure:"scheduledAt"`
	StartedAt    *time.Time `json:"startedAt,omitempty" mapstructure:"startedAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty" mapstructure:"completedAt"`
}

type DryRunEventData struct {
	EventData `mapstructure:",squash"`
	Status    string        `json:"status,omitempty" mapstructure:"status"`
	Errors    []DryRunError `json:"errors,omitempty" mapstructure:"errors"`
}

type DeploymentEventData struct {
	EventData     `mapstructure:",squash"`
	StageId       *uuid.UUID `json:"stageId,omitempty" mapstructure:"stageId"`
	CorrelationId *uuid.UUID `json:"correlationId,omitempty" mapstructure:"correlationId"`
	Message       string     `json:"message,omitempty" mapstructure:"message"`
}

type StageEventData struct {
	EventData         `mapstructure:",squash"`
	ParentOperationId *uuid.UUID `json:"parentOperationId,omitempty" mapstructure:"parentOperationId"`
	StageId           *uuid.UUID `json:"stageId,omitempty" mapstructure:"stageId"`
	CorrelationId     *uuid.UUID `json:"correlationId,omitempty" mapstructure:"correlationId"`
	Message           string     `json:"message,omitempty" mapstructure:"message"`
}

func (m *EventHookMessage) HashCode() string {
	hash32 := fnv.New32a()

	values := []string{
		m.HookId.String(),
		m.Type,
		m.Status,
		m.Subject,
	}

	if m.Data != nil {
		eventData := &EventData{}
		err := internal.Decode(m.Data, eventData)
		if err == nil && eventData != nil {
			values = append(values, eventData.OperationId.String())
		}
	}

	value := strings.Join(values, "")
	hash32.Write([]byte(value))

	hash := fmt.Sprintf("%x", hash32.Sum(nil))
	return hash
}

func (m *EventHookMessage) DryRunEventData() (*DryRunEventData, error) {
	return decode[DryRunEventData](EventDataTypePrefixDryRun, m)
}

func (m *EventHookMessage) StageEventData() (*StageEventData, error) {
	return decode[StageEventData](EventDataTypePrefixStage, m)
}

func (m *EventHookMessage) DeploymentEventData() (*DeploymentEventData, error) {
	return decode[DeploymentEventData](EventDataTypePrefixDeployment, m)
}

func (m *EventHookMessage) DeploymentId() (uint, error) {
	if data, ok := m.Data.(DeploymentEventData); ok {
		return uint(data.DeploymentId), nil
	}

	if m.Subject != "" && strings.HasPrefix(m.Subject, "/deployments/") {
		values := strings.Split(strings.TrimPrefix(m.Subject, "/"), "/")

		deploymentId, err := strconv.Atoi(values[1])
		if err != nil {
			return 0, err
		}
		return uint(deploymentId), nil
	}
	return 0, errors.New("unable to get deployment id using data or the subject")
}

func (m *EventHookMessage) SetSubject(deploymentId uint, stageId *uuid.UUID) {
	m.Subject = "/deployments/" + strconv.Itoa(int(deploymentId))
	if stageId != nil {
		m.Subject += "/stages/" + stageId.String()
	}
}

func decode[TData any](typePrefix string, m *EventHookMessage) (*TData, error) {
	if m == nil {
		return nil, errors.New("message is nil for type " + typePrefix)
	}

	if !strings.HasPrefix(m.Type, typePrefix) {
		return nil, errors.New("message event type is not " + typePrefix + "*")
	}

	var data *TData
	err := internal.Decode(m.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to decode data for type %T", *new(TData))
	}
	return data, nil
}
