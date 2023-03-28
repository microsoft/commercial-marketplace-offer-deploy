package test_test

import (
	//"context"
	//"log"
	"testing"
	// "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	// "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	// "github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/messaging"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/require"
)

type serviceBusSuite struct {
	suite.Suite
	ns string
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

func (s *serviceBusSuite) TestMessageSendSuccess() {
	sbConfig := messaging.ServiceBusConfig{
		Namespace: s.ns,
		QueueName: s.queueName,
	}
	config := messaging.PublisherConfig {
		Type: "servicebus",
		TypeConfig: sbConfig,
	}
	publisher, err := messaging.CreatePublisher(config)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), publisher)
	message := messaging.DeploymentMessage {
		Header: messaging.DeploymentMessageHeader {
			Topic: "testtopic",
		},
		Body: "testbody",
	}
	err = publisher.Publish(message)
	require.NoError(s.T(), err)
}