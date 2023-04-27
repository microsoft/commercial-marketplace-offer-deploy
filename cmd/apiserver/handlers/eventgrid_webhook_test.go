package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	testutils "github.com/microsoft/commercial-marketplace-offer-deploy/test/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var appConfig = getAppConfig()

func TestEventGridWebHook(t *testing.T) {
	//Setup
	json, err := testutils.ReaderFromJsonFile("testdata/eventgridevent.sb.json")
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/", json)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	// act & assert
	handlerFunc := getHandler(t)
	assert.NotNil(t, handlerFunc)

	err = handlerFunc(c)
	assert.NoError(t, err)

	log.Print(recorder.Body.String())
	assert.Equal(t, http.StatusOK, recorder.Code)
}

// setup handler with dependencies we control & fake
func getHandler(t *testing.T) echo.HandlerFunc {
	return func(c echo.Context) error {
		credential, err := azidentity.NewDefaultAzureCredential(nil)

		// we MUST have a credential or we're dead in the water anyway
		if err != nil {
			t.Fatalf("failed to create credential: %v", err)
		}

		db := setupDatabase()

		sender, err := newMessageSender(appConfig, credential)
		if err != nil {
			t.Fatal(err)
		}

		messageFactory, err := newWebHookEventMessageFactory(appConfig.Azure.SubscriptionId, db, credential)
		if err != nil {
			t.Fatal(err)
		}

		handler := eventGridWebHook{
			db:             db,
			messageFactory: messageFactory,
			sender:         sender,
		}

		log.Print("handler: ", handler)
		return handler.Handle(c)
	}
}

func getAppConfig() *config.AppConfig {
	appConfig := &config.AppConfig{}
	name := "test"
	config.LoadConfiguration("testdata", &name, appConfig)

	log.Printf("appConfig: %+v", appConfig)
	return appConfig
}

func setupDatabase() *gorm.DB {
	databaseOptions := &data.DatabaseOptions{
		UseInMemory: true,
	}
	db := data.NewDatabase(databaseOptions).Instance()
	db.Save(&data.Deployment{
		Name: "test-deployment",
	})
	return db
}
