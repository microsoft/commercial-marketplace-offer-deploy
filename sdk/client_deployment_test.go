package sdk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartDeployment(t *testing.T) {

	// Setup
	templateParameters := getParameters(t, "../test/testdata/nameviolation/success/")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, _ := io.ReadAll(r.Body)
		var received = make(map[string]interface{})
		json.Unmarshal(body, &received)

		if r.Method == "POST" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id": "test-id"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		//log.Printf("request: %+v", r)
		assert.Equal(t, "1", strings.Split(r.RequestURI, "/")[2])

		//get something from templateparams and assert that it is in the request body
		equals := reflect.DeepEqual(templateParameters, received)
		assert.True(t, equals)
		//create a utilty function to compare the two maps

	}))
	defer ts.Close()

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	client, err := NewClient(ts.URL, cred, nil)

	require.NoError(t, err)
	require.NotNil(t, client)

	var ctx context.Context = context.Background()

	// TODO: properly construct the startdeployment params
	// create
	_, err = client.Start(ctx, 1, templateParameters, nil)

	// Assertions
	if err != nil {
		t.Logf("Error: %s", err)
	}
}

func getParameters(t *testing.T, path string) map[string]interface{} {
	paramsPath := filepath.Join(path, "parameters.json")
	parameters, err := readJson(paramsPath)
	require.NoError(t, err)
	return parameters
}

func readJson(path string) (map[string]interface{}, error) {
	templateFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	template := make(map[string]interface{})
	if err := json.Unmarshal(templateFile, &template); err != nil {
		return nil, err
	}
	return template, nil
}
