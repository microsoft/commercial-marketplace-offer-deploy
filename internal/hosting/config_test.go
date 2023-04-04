package hosting

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
)

type TestConfigSection struct {
	Entry string `mapstructure:"TEST_CONFIGSECTION_ENTRY"`
}

type TestConfig struct {
	Entry   string `mapstructure:"TEST_CONFIG_ENTRY"`
	Section TestConfigSection
}

func TestEnvironmentVariablesLoad(t *testing.T) {
	os.Clearenv()

	value := "testvalue"
	os.Setenv("TEST_CONFIG_ENTRY", value)

	config := &TestConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), config)

	assert.NoError(t, err)
	assert.Equal(t, value, config.Entry)
}

func TestFileValuesLoadAndOverrideByEnvVar(t *testing.T) {
	os.Clearenv()

	value := "envvalue"
	os.Setenv("TEST_CONFIG_ENTRY", value)

	config := &TestConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), config)
	assert.NoError(t, err)
	assert.Equal(t, value, config.Entry)
}

func TestFileValuesLoad(t *testing.T) {
	os.Clearenv()

	config := &TestConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), config)
	assert.NoError(t, err)
	assert.Equal(t, "filevalue", config.Entry)
}

func TestLoadWithValueMissingFromEnvFile(t *testing.T) {
	os.Clearenv()

	os.Setenv("TEST_CONFIGSECTION_ENTRY", "configsectionvalue")

	configSection := &TestConfigSection{}

	err := LoadConfiguration("./testdata", to.Ptr("test"), configSection)

	assert.NoError(t, err)
	assert.Equal(t, "configsectionvalue", configSection.Entry)
}

func TestConfigSectionsLoad(t *testing.T) {
	os.Clearenv()

	os.Setenv("TEST_CONFIGSECTION_ENTRY", "configsectionvalue")

	config := &TestConfig{
		Section: TestConfigSection{},
	}

	err := LoadConfiguration("./testdata", to.Ptr("test"), config)

	assert.NoError(t, err)
	assert.Equal(t, "configsectionvalue", config.Section.Entry)
}
