package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartDeployment(t *testing.T) {
	db := data.NewDatabase(&data.DatabaseOptions{Dsn: "./testdata/test.db"}).Instance()
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(deploymentJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := CreateDeployment(c, db) // create a deployment / http POST to the server
	if err != nil {
		t.Fatal(err)
	}
	var deploymentResult api.Deployment
	json.Unmarshal(rec.Body.Bytes(), &deploymentResult) // parse the response from the apiserver and map to an object: deploymentResult

	templateParameters := getParameters(t, "../../../test/testdata/nameviolation/success/")
	deploymentName := "StartDeploymentTest"

	invokeDeployOperation := api.InvokeDeploymentOperation{
		Parameters: templateParameters,
		Name:       &deploymentName,
	}
	StartDeployment(int(*deploymentResult.ID), invokeDeployOperation, db) // start the deployment / http POST to the server

	var startDeploymentResult api.InvokedOperation
	json.Unmarshal(rec.Body.Bytes(), &startDeploymentResult) // parse the response and map to the INVOKEDOperation object: startDeploymentResult

	retrieved := &data.Deployment{}
	db.First(&retrieved, *deploymentResult.ID) // query result matching the deployment ID; also mapped to object retrieved(.ID)

	log.Printf("value back from DB: %v", retrieved.ID)
	assert.Equal(t, *deploymentResult.ID, int32(retrieved.ID)) // validate the database saved the state
}
func getParameters(t *testing.T, path string) map[string]interface{} {
	paramsPath := filepath.Join(path, "parameters.json")
	parameters, err := utils.ReadJson(paramsPath)
	require.NoError(t, err)
	return parameters
}
