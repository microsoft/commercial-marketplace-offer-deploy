package hosting

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	entry string `mapstructure:"TEST_CONFIG_ENTRY"`
}

func TestLoadConfiguration(t *testing.T) {
	os.Setenv("TEST_CONFIG_ENTRY", "testentry")
	config := testConfig{}
	err := LoadConfiguration("./testdata", nil, &config)

	require.NoError(t, err)
	assert.Equal(t, "testentry", config.entry)
}
