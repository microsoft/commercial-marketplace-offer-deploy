package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	testutils "github.com/microsoft/commercial-marketplace-offer-deploy/test/utils"
	"github.com/stretchr/testify/assert"
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
)

func ToMove_TestCreateDeployment(t *testing.T) {
	// Setup
	deploymentJson := getFakeCreateDeploymentJson(t)

	db := data.NewDatabase(&data.DatabaseOptions{Dsn: "./testdata/test.db"}).Instance()
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(deploymentJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	c := e.NewContext(request, recorder)

	handler := createDeploymentHandler{
		db:     db,
		mapper: mapper.NewCreateDeploymentMapper(),
	}

	//execute handler
	err := handler.Handle(c)

	// Assertions
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// check response
	var jsonResponse api.Deployment
	json.Unmarshal(recorder.Body.Bytes(), &jsonResponse)
	assert.Equal(t, "test name", *jsonResponse.Name)

	// check database
	id := uint(*jsonResponse.ID)
	data := &data.Deployment{}
	db.First(data, id)

	assert.Equal(t, id, data.ID)
	assert.Equal(t, *jsonResponse.Name, data.Name)

	// make sure stages are created
	assert.Len(t, data.Stages, 1)
	assert.Equal(t, "Test Stage", data.Stages[0].Name)
}

func getFakeCreateDeploymentJson(t *testing.T) string {
	template, err := testutils.NewFromJsonFile[map[string]any]("testdata/azuredeploy.json")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	item := &api.CreateDeployment{
		Name:           to.Ptr("test name"),
		SubscriptionID: to.Ptr("test"),
		ResourceGroup:  to.Ptr("test"),
		Location:       to.Ptr("test"),
		Template:       template,
	}

	bytes, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	return string(bytes)
}
