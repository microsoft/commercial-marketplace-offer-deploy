package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfiguration(t *testing.T) {
	name := "test"
	config, err := LoadConfiguration("./testdata", &name)

	require.NoError(t, err)
	require.NotNil(t, config)

	require.Equal(t, "test-subscription-id", config.Azure.SubscriptionId)
}
