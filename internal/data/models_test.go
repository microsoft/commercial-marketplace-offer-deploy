package data

import (
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeploymentTemplateMarshaling(t *testing.T) {
	template := make(map[string]interface{})

	attr := make(map[string]interface{})
	attr["field"] = "test"
	template["mapField"] = attr

	model := &Deployment{Template: template}

	db := NewDatabase(&DatabaseOptions{UseInMemory: true}).Instance()
	db.Save(model)

	var result *Deployment
	db.First(&result)

	log.Print(result.Template)
	require.NotNil(t, result)
	require.Equal(t, "test", result.Template["mapField"].(map[string]interface{})["field"])
}

func TestInvokedOperationUpdate(t *testing.T) {
	params := make(map[string]interface{})
	params["stagedId"] = uuid.New()

	model := &InvokedOperation{
		DeploymentId: 1,
		Parameters:   params,
		Name:         string(operation.StatusScheduled),
		Status:       "test",
		Result:       make(map[string]interface{}),
	}

	// save
	db := NewDatabase(&DatabaseOptions{UseInMemory: true}).Instance()
	db.Save(model)

	// read
	var result *InvokedOperation
	db.First(&result, model.ID)

	// update
	result.Status = "updated"
	result.Result = nil
	db.Save(result)

	assert.NotEqualValues(t, uuid.Nil, result.ID)
}

func TestDeploymentUpdate(t *testing.T) {
	params := make(map[string]interface{})
	params["stagedId"] = uuid.New()

	model := &Deployment{
		Name:     "test",
		Template: params,
	}

	// save
	db := NewDatabase(&DatabaseOptions{UseInMemory: true}).Instance()
	db.Save(model)

	// read
	var result *Deployment
	db.First(&result, model.ID)

	// update
	result.Name = result.Name + "updated"
	db.Save(result)
}
