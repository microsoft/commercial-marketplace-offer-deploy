package handlers

import (
	"errors"
	"log"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
<<<<<<< HEAD
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"gorm.io/gorm"
)

func StartDeployment(deploymentId int, operation api.InvokeDeploymentOperation, db *gorm.DB) (interface{}, error) {
	log.Printf("Inisde CreateDryRun deploymentId: %d", deploymentId)
=======
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/messaging"
	"gorm.io/gorm"
)

func StartDeployment(deploymentId int, operation generated.InvokeDeploymentOperation, db *gorm.DB) (interface{}, error) {
	log.Printf("Inside StartDeployment")
>>>>>>> main

	//gather data: deploymentId
	retrieved := &data.Deployment{}
	db.First(&retrieved, deploymentId)

	templateParams := operation.Parameters
	if templateParams == nil {
		return nil, errors.New("templateParams were not provided")
	}

	//do the work to update the database and post a message to the service bus

	// TODO: WRAP in transaction

	//// deployment.Pending
	config := messaging.PublisherConfig{}
	config.Type = "servicebus"
	publisher, _ := messaging.CreatePublisher(config)

	message := messaging.DeploymentMessage{
		Header: messaging.DeploymentMessageHeader{
			Topic: "TestTopic",
		},
		Body: "TestContent",
	}
	err := publisher.Publish(message)

	// Update DB

	//End transaction

	if err != nil {
		return nil, err
	}

	// formulate the response
	timestamp := time.Now().UTC()
	status := "OK"

<<<<<<< HEAD
	// res := deployment.DryRun(&azureDeployment)
	// uuid := uuid.New().String()
	// timestamp := time.Now().UTC()
	// status := "OK"
	returnedResult := api.InvokedOperation{
		ID:        nil,
		InvokedOn: nil,
		Result:    nil,
		Status:    nil,
=======
	returnedResult := generated.InvokedOperation{
		ID:         operation.Name,
		Parameters: operation.Parameters.(map[string]interface{}),
		InvokedOn:  &timestamp,
		Result:     nil,
		Status:     &status,
>>>>>>> main
	}
	return returnedResult, nil
}
