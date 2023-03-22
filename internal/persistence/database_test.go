package persistence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDatabaseWithNilOptions(t *testing.T) {
	DefaultDatabasePath = "./testdata"

	db := NewDatabase(nil)
	require.NotNil(t, db)
}
