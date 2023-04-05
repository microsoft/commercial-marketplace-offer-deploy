package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	testutils "github.com/microsoft/commercial-marketplace-offer-deploy/test/utils"
	"github.com/stretchr/testify/assert"
)

func TestEventGridWebHook(t *testing.T) {
	// Setup
	json, err := testutils.ReaderFromJsonFile("testdata/eventgrid.json")
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/", json)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	// act
	err = EventGridWebHook(c, nil, nil, nil)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
