package mapper

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	testutils "github.com/microsoft/commercial-marketplace-offer-deploy/test/utils"
	"github.com/stretchr/testify/assert"
)

// create a test that creates a new mapper then takes in a CreateDeployment and test the mapped output
func TestCreateDeploymentMapper(t *testing.T) {
	input := getFakeCreateDeployment(t)

	assert.NotNil(t, input)
	assert.NotNil(t, input.Template)

	mapper := NewCreateDeploymentMapper()
	result, err := mapper.Map(input)

	assert.NoError(t, err)
	assert.EqualValues(t, *input.Name, result.Name)

	assert.Len(t, result.Stages, 1)
	assert.EqualValues(t, "Test Stage", result.Stages[0].Name)
}

func getFakeCreateDeployment(t *testing.T) *sdk.CreateDeployment {
	template, err := testutils.NewFromJsonFile[map[string]any]("testdata/azuredeploy.json")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	item := &sdk.CreateDeployment{
		Name:           to.Ptr("test name"),
		SubscriptionID: to.Ptr("test"),
		ResourceGroup:  to.Ptr("test"),
		Location:       to.Ptr("test"),
		Template:       template,
	}
	return item
}
