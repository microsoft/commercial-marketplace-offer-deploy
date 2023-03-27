package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDirectoryPath string = "./testdata"

func TestNewDatabaseWithNilOptions(t *testing.T) {
	setupTestDirectory(t)

	// override default db path to /testdata
	SetDefaultDatabasePath(testDirectoryPath)

	db := NewDatabase(nil)
	require.NotNil(t, db)

	// assert db file exists
	dbFileCreatedByNewDatabase := filepath.Join(testDirectoryPath, DatabaseFileName)

	_, err := os.Stat(dbFileCreatedByNewDatabase)
	require.NoError(t, err)

	cleanTestDirectory()
}

func setupTestDirectory(t *testing.T) {
	if _, err := os.Stat(testDirectoryPath); err != nil {
		err := os.Mkdir(testDirectoryPath, 0755)
		require.NoError(t, err)
	}
}

func cleanTestDirectory() {
	//clean up
	os.RemoveAll(testDirectoryPath)
}
