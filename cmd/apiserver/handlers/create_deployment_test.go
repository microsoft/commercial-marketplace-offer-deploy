package handlers

import (
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/stretchr/testify/require"
)

func TestSaveDeployment(t *testing.T) {
	data.DefaultDatabasePath = "./testdata"
	db := data.NewDatabase(nil).Instance()
	require.NotNil(t, db)
}
