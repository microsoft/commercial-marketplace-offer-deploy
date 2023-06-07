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

func (m *EventHookMessage) GetHash() string {
	hash := fnv.New32a()
	hashString := fmt.Sprintf("%s%s%s%s", m.HookId, m.Type, m.Status, m.Subject)
	hash.Write([]byte(hashString))
	hashBytes := hash.Sum(nil)
	hashValue := fmt.Sprintf("%x", hashBytes)
	return hashValue
}

func (m *EventHookMessage) DryRunEventData() (*DryRunEventData, error) {
	if m.Type != EventTypeDryRunCompleted.String() {
		return nil, errors.New("message type is not dryRunCompleted")
	}

	var data *DryRunEventData
	err := internal.Decode(m.Data, &data)
	if err != nil {
		return nil, errors.New("data is not of type DryRunEventData")
	}
	return data, nil
}

func (m *EventHookMessage) DeploymentEventData() (*DeploymentEventData, error) {
	if strings.HasPrefix(m.Type, "deployment") {
		var data *DeploymentEventData
		err := internal.Decode(m.Data, &data)
		if err != nil {
			return nil, errors.New("data is not of type DeploymentEventData")
		}
		return data, nil
	}
	return nil, errors.New("message event type is not deployment*")
}

// Event data for a message

type DryRunEventData struct {
	DeploymentId int           `json:"deploymentId" mapstructure:"deploymentId"`
	OperationId  uuid.UUID     `json:"operationId" mapstructure:"operationId"`
	Attempts     int           `json:"attempts" mapstructure:"attempts"`
	Status       string        `json:"status,omitempty" mapstructure:"status"`
	Errors       []DryRunError `json:"errors,omitempty" mapstructure:"errors"`
	StartedAt    time.Time     `json:"startedAt,omitempty" mapstructure:"startedAt"`
	CompletedAt  time.Time     `json:"completedAt,omitempty" mapstructure:"completedAt"`
}

type DeploymentEventData struct {
	DeploymentId int        `json:"deploymentId,omitempty" mapstructure:"deploymentId"`
	StageId      *uuid.UUID `json:"stageId,omitempty" mapstructure:"stageId"`
	OperationId  uuid.UUID  `json:"operationId,omitempty" mapstructure:"operationId"`

	// the correlation ID used to track azure deployments. This may be nil if the particular event data is about something that
	// happened inside MODM and not azure.
	CorrelationId *uuid.UUID `json:"correlationId,omitempty" mapstructure:"correlationId"`
	Attempts      int        `json:"attempts,omitempty" mapstructure:"attempts"`
	Message       string     `json:"message,omitempty" mapstructure:"message"`
	StartedAt     time.Time  `json:"startedAt,omitempty" mapstructure:"startedAt"`
	CompletedAt   time.Time  `json:"completedAt,omitempty" mapstructure:"completedAt"`
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
