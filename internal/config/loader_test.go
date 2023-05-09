package config

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
)

type TestConfigSection struct {
	Entry string `mapstructure:"TEST_CONFIGSECTION_ENTRY"`
}

type TestConfig struct {
	Entry      string   `mapstructure:"TEST_CONFIG_ENTRY"`
	ArrayEntry []string `mapstructure:"TEST_CONFIG_ARRAY_ENTRY"`
	Section    TestConfigSection
}

func TestEnvironmentVariablesLoadAndOverrideEnvFile(t *testing.T) {
	os.Clearenv()

	value := "envvalue"
	os.Setenv(EnvironmentVariablePrefix+"_TEST_CONFIG_ENTRY", value)

	config := &TestConfig{}
	err := LoadConfiguration("./testdata", to.Ptr("test"), config)

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
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

	os.Setenv(EnvironmentVariablePrefix+"_TEST_CONFIGSECTION_ENTRY", "configsectionvalue")

	configSection := &TestConfigSection{}

	err := LoadConfiguration("./testdata", to.Ptr("test"), configSection)

	assert.NoError(t, err)
	assert.Equal(t, "configsectionvalue", configSection.Entry)
}

func TestConfigSectionsLoad(t *testing.T) {
	os.Clearenv()

	os.Setenv(EnvironmentVariablePrefix+"_TEST_CONFIGSECTION_ENTRY", "configsectionvalue")

	config := &TestConfig{
		Section: TestConfigSection{},
	}

	err := LoadConfiguration("./testdata", to.Ptr("test"), config)

	assert.NoError(t, err)
	assert.Equal(t, "configsectionvalue", config.Section.Entry)
}

func TestConfigArrayValue(t *testing.T) {
	os.Clearenv()

	os.Setenv(EnvironmentVariablePrefix+"_TEST_CONFIG_ARRAY_ENTRY", "item1,item2,item3")

	config := &TestConfig{
		Section: TestConfigSection{},
	}

	err := LoadConfiguration("./testdata", to.Ptr("test"), config)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(config.ArrayEntry))
	assert.Equal(t, "item1", config.ArrayEntry[0])
}

func TestEnvPrefix(t *testing.T) {
	os.Clearenv()
	os.Setenv(EnvironmentVariablePrefix+"_TEST_CONFIG_ARRAY_ENTRY", "envvarprefix1")

	config := &TestConfig{
		Section: TestConfigSection{},
	}

	err := LoadConfiguration("./testdata", to.Ptr("test"), config)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(config.ArrayEntry))
	assert.Equal(t, "envvarprefix1", config.ArrayEntry[0])
}
