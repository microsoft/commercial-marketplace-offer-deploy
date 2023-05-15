package config

import (
	"net/url"
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
	url, _ := url.Parse("https://test.com")
	appConfig := &AppConfig{
		Http: HttpSettings{
			BaseUrl: url.String(),
		},
	}
	assert.Equal(t, url, appConfig.GetPublicBaseUrl())

	appConfig.Http.BaseUrl = "https://test.com/"
	url, _ = url.Parse(appConfig.Http.BaseUrl)

	assert.Equal(t, url, appConfig.GetPublicBaseUrl())
}
