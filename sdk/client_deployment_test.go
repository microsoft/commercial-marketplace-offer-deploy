package sdk

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/handlers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	db             = data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	deploymentJson = `{
		"name":"test-deployment", 
		"subscriptionId":"test-id",
		"resourceGroup":"test-rg",
		"location":"testus",
		"template": {}
	}`
	templateParameters map[string]interface{}
)

const testEndPoint = "http://localhost:8080"

func TestStartDeployment(t *testing.T) {
	// Setup
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	client, err := NewClient(testEndPoint, cred, nil)

	require.NoError(t, err)
	require.NotNil(t, client)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(deploymentJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.CreateDeployment(c, db)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	var results generated.Deployment
	json.Unmarshal(rec.Body.Bytes(), &results)

	// TODO: properly construct the startdeployment params
	result, err := client.StartDeployment(context.Background(), *results.ID, templateParameters)

	// Assertions
	if err != nil {
		t.Logf("Error: %s", err)
	}
	require.NotNil(t, result)
}
