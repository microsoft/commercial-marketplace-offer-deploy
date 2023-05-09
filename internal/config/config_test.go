package config

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
)

func TestAppConfigWithEnvPrefix(t *testing.T) {
	logfilePath := "/test/path"
	os.Setenv("MODM_LOG_FILE_PATH", logfilePath)

	appConfig := &AppConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), appConfig)
	assert.NoError(t, err)
	assert.Equal(t, logfilePath, appConfig.Logging.FilePath)
}
