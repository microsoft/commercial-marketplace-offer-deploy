package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
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

func TestCreateDeployment(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(deploymentJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, CreateDeployment(c, db)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var result api.Deployment
		json.Unmarshal(rec.Body.Bytes(), &result)
		assert.Equal(t, int32(1), *result.ID)
		assert.Equal(t, "test-deployment", *result.Name)
	}
}
