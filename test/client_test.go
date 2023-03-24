package test_test

import (
	"context"
	"log"
	"testing"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/require"
)

type clientSuite struct {
	suite.Suite
	endpoint string
}

func (s *clientSuite) SetupSuite() {
	s.endpoint = "http://localhost:8080"
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, &clientSuite{})
}

func (s *clientSuite) TestNewClient() {
	client, err := sdk.NewClient("test", nil, nil)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), client)

	cloudConfig := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "test", Services: map[cloud.ServiceName]cloud.ServiceConfiguration{},
	}
	_, err = sdk.NewClient("test", nil, &sdk.ClientOptions{ClientOptions: azcore.ClientOptions{Cloud: cloudConfig}})
	require.NoErrorf(s.T(), err, "No error.")
}

func (s *clientSuite) TestListDeployments() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	client, err := sdk.NewClient(s.endpoint, cred, nil)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), client)

	result, err := client.ListDeployments(context.TODO())

	if err != nil {
		s.T().Logf("Error: %s", err)
	}
	require.NotNil(s.T(), result)
}
