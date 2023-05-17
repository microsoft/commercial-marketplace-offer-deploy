package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/assert"
)

func TestGetDeployment(t *testing.T) {
	//new in-memory db
	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()

	//Add deployment to db
	data := &data.Deployment{}
	db.Save(data)

	//Setup
	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	c.SetPath("/:deploymentId")
	c.SetParamNames("deploymentId")
	c.SetParamValues(strconv.Itoa(1))

	handler := getDeploymentHandler{
		db:     db,
		mapper: mapper.NewDeploymentMapper(),
	}

	//execute handler
	err := handler.Handle(c)
	assert.NoError(t, err)

	//check for invalid ID
	id, err := handler.getDeploymentId(c)
	if err != nil {
		t.Error("Invalid ID")
	}

	//deployment
	var jsonResponse sdk.Deployment
	json.Unmarshal(recorder.Body.Bytes(), &jsonResponse)

	assert.Equal(t, data.Name, *jsonResponse.Name)
	responseID := uint(*jsonResponse.ID)

	//check if given ID matches db ID
	assert.EqualValues(t, id, responseID)
}
