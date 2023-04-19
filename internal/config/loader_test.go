package config

import (
	"log"
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

func ToMove_TestEnvironmentVariablesLoad(t *testing.T) {
	os.Clearenv()

	value := "testvalue"
	os.Setenv("TEST_CONFIG_ENTRY", value)

	config := &TestConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), config)

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, value, config.Entry)
}

func ToMove_TestFileValuesLoadAndOverrideByEnvVar(t *testing.T) {
	os.Clearenv()

	value := "envvalue"
	os.Setenv("TEST_CONFIG_ENTRY", value)

	config := &TestConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), config)
	assert.NoError(t, err)
	assert.Equal(t, value, config.Entry)
}

func ToMove_TestFileValuesLoad(t *testing.T) {
	os.Clearenv()

	config := &TestConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), config)
	assert.NoError(t, err)
	assert.Equal(t, "filevalue", config.Entry)
}

func ToMove_TestLoadWithValueMissingFromEnvFile(t *testing.T) {
	os.Clearenv()

	os.Setenv("TEST_CONFIGSECTION_ENTRY", "configsectionvalue")

	configSection := &TestConfigSection{}

	err := LoadConfiguration("./testdata", to.Ptr("test"), configSection)

	assert.NoError(t, err)
	assert.Equal(t, "configsectionvalue", configSection.Entry)
}

func ToMove_TestConfigSectionsLoad(t *testing.T) {
	os.Clearenv()

	os.Setenv("TEST_CONFIGSECTION_ENTRY", "configsectionvalue")

	config := &TestConfig{
		Section: TestConfigSection{},
	}

	err := LoadConfiguration("./testdata", to.Ptr("test"), config)

	assert.NoError(t, err)
	assert.Equal(t, "configsectionvalue", config.Section.Entry)
}
