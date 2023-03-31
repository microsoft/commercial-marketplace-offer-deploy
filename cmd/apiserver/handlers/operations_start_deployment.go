package handlers

import (
	"errors"
	"log"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
	"gorm.io/gorm"
)

func StartDeployment(deploymentId int, operation generated.InvokeDeploymentOperation, db *gorm.DB) (interface{}, error) {
	log.Printf("Inside StartDeployment")

	//gather data: deploymentId
	retrieved := &data.Deployment{}
	db.First(&retrieved, deploymentId)

	templateParams := operation.Parameters
	if templateParams == nil {
		return nil, errors.New("templateParams were not provided")
	}

	//do the work to update the database and post a message to the service bus
	//// deployment.Pending

	// formulate the response
	timestamp := time.Now().UTC()
	status := "OK"

	returnedResult := generated.InvokedOperation{
		ID:         operation.Name,
		Parameters: operation.Parameters.(map[string]interface{}),
		InvokedOn:  &timestamp,
		Result:     nil,
		Status:     &status,
	}
	return returnedResult, nil
}
