package config

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
)

func TestAppConfigWithEnvPrefix(t *testing.T) {
	logfilePath := "/test/path/to/logfile"
	dbPath := "/test/path/to/db"

	os.Setenv("MODM_DB_PATH", dbPath)
	os.Setenv("MODM_LOG_FILE_PATH", logfilePath)

	appConfig := &AppConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), appConfig)
	assert.NoError(t, err)

	assert.Equal(t, logfilePath, appConfig.Logging.FilePath)
	assert.Equal(t, dbPath, appConfig.Database.Path)
}

func Test_AppConfig_GetPublicBaseUrl(t *testing.T) {
	// without trailing slash
	appConfig := &AppConfig{
		Http: HttpSettings{
			BaseUrl: "https://test.com",
		},
	}
	assert.Equal(t, "https://test.com/", appConfig.GetPublicBaseUrl())

	appConfig.Http.BaseUrl = "https://test.com/"
	assert.Equal(t, appConfig.Http.BaseUrl, appConfig.GetPublicBaseUrl())
}
