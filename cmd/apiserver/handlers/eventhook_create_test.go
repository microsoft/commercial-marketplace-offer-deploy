package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	testutils "github.com/microsoft/commercial-marketplace-offer-deploy/test/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateEventHookIsIdempotentByName(t *testing.T) {
	db := getInMemoryDb()

	json, err := testutils.ReaderFromJsonFile("testdata/eventhook.json")
	assert.NoError(t, err)

	e := echo.New()
	defer e.Close()

	request := httptest.NewRequest(http.MethodPost, "/", json)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	//act
	err = CreateEventHook(c, db)
	assert.NoError(t, err)

	t.Log(recorder.Body.String())
	assert.Equal(t, http.StatusOK, recorder.Code)

	json, _ = testutils.ReaderFromJsonFile("testdata/eventhook.json")
	request = httptest.NewRequest(http.MethodPost, "/", json)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder = httptest.NewRecorder()
	c = e.NewContext(request, recorder)

	err = CreateEventHook(c, db)
	assert.NoError(t, err)
	t.Log(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func getInMemoryDb() *gorm.DB {
	databaseOptions := &data.DatabaseOptions{
		UseInMemory: true,
	}
	db := data.NewDatabase(databaseOptions).Instance()
	return db
}
