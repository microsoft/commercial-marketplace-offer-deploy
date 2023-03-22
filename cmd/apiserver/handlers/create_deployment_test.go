package handlers

import (
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/generated"
	"github.com/stretchr/testify/require"
)

func TestSaveDeployment(t *testing.T) {
	data.DefaultDatabasePath = "./testdata"
	db := data.NewDatabase(nil).Instance()

	name := "test"
	command := generated.CreateDeployment{
		Name:     &name,
		Template: "",
	}
	result, err := saveDeployment(command, db)

	require.NotNil(t, result)
	require.NoError(t, err)
}
