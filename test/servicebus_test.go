package test_test

import (
	//"context"
	//"log"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	// "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	// "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	// "github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/messaging"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type serviceBusSuite struct {
	suite.Suite
	ns        string
	queueName string
}

func TestServiceBusSuite(t *testing.T) {
	suite.Run(t, &serviceBusSuite{})
}

func (s *serviceBusSuite) SetupSuite() {
	s.ns = "bobjacmodm.servicebus.windows.net"
	s.queueName = "deployeventqueue"
}

func (s *serviceBusSuite) SetupTest() {
	// create the service bus namespace with the ns and queuename
}

func (s *serviceBusSuite) publishTestMessage(topicHeader string, body string) {
	sbConfig := messaging.ServiceBusConfig{
		Namespace: s.ns,
		QueueName: s.queueName,
	}
	config := messaging.PublisherConfig{
		Type:       "servicebus",
		TypeConfig: sbConfig,
	}
	publisher, err := messaging.CreatePublisher(config)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), publisher)
	message := messaging.DeploymentMessage{
		Header: messaging.DeploymentMessageHeader{
			Topic: topicHeader,
		},
		Body: body,
	}
	err = publisher.Publish(message)
	require.NoError(s.T(), err)
}

func (s *serviceBusSuite) TestMessageSendSuccess() {
	for i := 0; i < 15; i++ {
		body := fmt.Sprintf("testbody%d", i)
		s.publishTestMessage("testtopic", body)
	}
}

func (s *serviceBusSuite) TestMessageReceiveSuccess() {
	sbConfig := messaging.ServiceBusConfig{
		Namespace: s.ns,
		QueueName: s.queueName,
	}

	handler := &fakeHandler{}

	receiver, err := messaging.NewServiceBusReceiver(sbConfig, handler)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receiver)
	fmt.Println("calling start")
	go receiver.Start()
	// sleep for 5 seconds to allow the receiver to start
	fmt.Println("starting sleep 1")
	time.Sleep(5 * time.Second)
	go receiver.Stop()
	fmt.Println("After the stop in TestMessageReceiveSuccess")
	fmt.Println("Starting sleep 2")
	time.Sleep(5 * time.Second)
	fmt.Println("After the second sleep")
}

type fakeHandler struct {
}

func (h *fakeHandler) Handle(ctx context.Context, message *azservicebus.ReceivedMessage) error {
	log.Println("Handling message")
	return nil
}
