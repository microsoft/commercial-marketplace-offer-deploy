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
	Entry   string `mapstructure:"TEST_CONFIG_ENTRY"`
	Section TestConfigSection
}

func TestMain(m *testing.M) {
	path := "./testdata"
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return
		}
		file.Close()
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	_, err = file.WriteString("TEST_CONFIG_ENTRY=filevalue \n")
}

func isError(err error) bool {
	return (err != nil)
}

func TestEnvironmentVariablesLoad(t *testing.T) {
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
	// needs test.env
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
