package test_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type serviceBusSuite struct {
	suite.Suite
	appConfig           *config.AppConfig
	testDirectory       string
	ns                  string
	eventsQueueName     string
	operationsQueueName string
	subscriptionId      string
	resourceGroupName   string
	location            string
	deploymentName      string
	deploymentId        uint
	invokedOperationId  uuid.UUID
	db                  data.Database
}

func TestServiceBusSuite(t *testing.T) {
	suite.Run(t, &serviceBusSuite{})
}

func (s *serviceBusSuite) SetupSuite() {
	testDirectory := "./testdata"

	//load config from testdata/.env
	appConfig := &config.AppConfig{}
	config.LoadConfiguration(testDirectory, nil, appConfig)

	log.Printf("SetupSuite - appConfig %v", appConfig)
	s.appConfig = appConfig

	s.ns = appConfig.Azure.GetFullQualifiedNamespace()
	s.eventsQueueName = string(messaging.QueueNameEvents)
	s.operationsQueueName = string(messaging.QueueNameOperations)
	s.subscriptionId = appConfig.Azure.SubscriptionId
	s.resourceGroupName = appConfig.Azure.ResourceGroupName
	s.location = appConfig.Azure.Location
	s.deploymentName = "test-deployment"

	s.testDirectory = testDirectory
	s.setupTestDirectory()
	data.SetDefaultDatabasePath(s.testDirectory)

	s.db = data.NewDatabase(nil)
}

func (s *serviceBusSuite) SetupTest() {
	s.createDeploymentForTests()
}

func (s *serviceBusSuite) setupTestDirectory() {
	if _, err := os.Stat(s.testDirectory); err != nil {
		err := os.Mkdir(s.testDirectory, 0755)
		require.NoError(s.T(), err)
	}
}

func (s *serviceBusSuite) createDeploymentForTests() {
	testDeploymentPath := "testdata/taggeddeployment"
	fullPath := filepath.Join(testDeploymentPath, "mainTemplateBicep.json")
	template, err := utils.ReadJson(fullPath)
	require.NoError(s.T(), err)

	paramsPath := filepath.Join(testDeploymentPath, "parametersBicep.json")
	parameters, err := utils.ReadJson(paramsPath)
	require.NoError(s.T(), err)

	deployment := &data.Deployment{
		Name:           s.deploymentName,
		Status:         "New",
		SubscriptionId: s.subscriptionId,
		ResourceGroup:  s.resourceGroupName,
		Location:       s.location,
		Template:       template,
	}

	s.db.Instance().Create(deployment)
	s.deploymentId = deployment.ID

	invokedOperation := &data.InvokedOperation{
		DeploymentId: deployment.ID,
		Parameters:   parameters,
	}

	s.db.Instance().Create(invokedOperation)
	s.invokedOperationId = invokedOperation.ID
}

func (s *serviceBusSuite) TestMessageSendSuccess() {
	log.Print("TestMessageSendSuccess")
	for i := 0; i < 15; i++ {
		text := fmt.Sprintf("testbody%d", i)
		message := &testMessage{Id: uuid.New().String(), Text: text}
		s.sendTestMessage(s.eventsQueueName, message)
	}
}

func (s *serviceBusSuite) TestOperationsSendSuccess() {
	var dataDeployment data.Deployment
	s.db.Instance().First(&dataDeployment, s.deploymentId)

	var invokedOperation data.InvokedOperation
	s.db.Instance().First(&invokedOperation, s.invokedOperationId)

	message := &messaging.InvokedOperationMessage{
		OperationId: invokedOperation.ID.String(),
	}
	s.sendTestMessage(s.operationsQueueName, message)
}

func (s *serviceBusSuite) TestMessageReceiveSuccess() {
	receiver := s.getMessageReceiver(s.eventsQueueName, s.getTestHandler())

	fmt.Println("calling start")
	go receiver.Start()

	//sleep for 5 seconds to allow the receiver to start
	fmt.Println("starting sleep 1")
	time.Sleep(5 * time.Second)

	go receiver.Stop()
	fmt.Println("After the stop in TestMessageReceiveSuccess")
	fmt.Println("Starting sleep 2")
	time.Sleep(5 * time.Second)
	fmt.Println("After the second sleep")
}

func (s *serviceBusSuite) TestOperationDeploymentSuccess() {
	handler := &testOperationsHandler{}
	receiver := s.getMessageReceiver(s.operationsQueueName, handler)

	fmt.Println("calling start")
	go receiver.Start()

	time.Sleep(60 * time.Minute)
	go receiver.Stop()
}

// helpers

func (s *serviceBusSuite) sendTestMessage(queueName string, message any) {
	ctx := context.TODO()
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(s.T(), err)

	options := messaging.MessageSenderOptions{
		SubscriptionId:          s.subscriptionId,
		Location:                s.location,
		ResourceGroupName:       s.resourceGroupName,
		FullyQualifiedNamespace: s.ns,
	}
	sender, err := messaging.NewServiceBusMessageSender(cred, options)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), sender)

	log.Printf("sending message to queue: %s", queueName)
	_, err = sender.Send(ctx, queueName, message)
	require.NoError(s.T(), err)
}

func (s *serviceBusSuite) getMessageReceiver(queueName string, handler any) messaging.MessageReceiver {
	options := s.getReceiverOptions(queueName)
	receiver, err := messaging.NewServiceBusReceiver(handler, options)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), receiver)

	return receiver
}

func (s *serviceBusSuite) getReceiverOptions(queueName string) messaging.ServiceBusMessageReceiverOptions {
	options := messaging.ServiceBusMessageReceiverOptions{
		MessageReceiverOptions:  messaging.MessageReceiverOptions{QueueName: queueName},
		FullyQualifiedNamespace: s.appConfig.Azure.GetFullQualifiedNamespace(),
	}
	return options
}

func (s *serviceBusSuite) getTestHandler() *testHandler {
	return &testHandler{}
}

type testMessage struct {
	Id   string
	Text string
}

type testHandler struct {
}

func (h *testHandler) Handle(message *testMessage, context messaging.MessageHandlerContext) error {
	log.Printf("Handling message [%s] - %s", message.Id, message.Text)
	return nil
}

type testOperationsHandler struct {
}

func (h *testOperationsHandler) Handle(message *messaging.InvokedOperationMessage, context messaging.MessageHandlerContext) error {
	log.Printf("Handling invoked operation message [%s]", message.OperationId)
	return nil
}
