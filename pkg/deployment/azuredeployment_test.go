package deployment_test

import (
	"path/filepath"
	"testing"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/stretchr/testify/assert"
)

func TestStartDeployment(t *testing.T) {
	fullPath := filepath.Join("../../test/testdata/nameviolation/nestedfailure", "mainTemplate.json")
	template, err := utils.ReadJson(fullPath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	resources := deployment.FindResourcesByType(template, "Microsoft.Resources/deployments")
	assert.Greater(t, len(resources), 0)
}