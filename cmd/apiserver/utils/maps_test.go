package utils

import (
	"encoding/json"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabaseWithNilOptions(t *testing.T) {
	var template map[string]any
	
	deploymentJsonString := "{\"$schema\":\"https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#\",\"contentVersion\":\"1.0.0.0\",\"metadata\":{\"_generator\":{\"name\":\"bicep\",\"templateHash\":\"1109253174677965668\",\"version\":\"0.15.31.15270\"}},\"parameters\":{\"aapName\":{\"defaultValue\":\"[format('aap{0}', uniqueString(utcNow(), resourceGroup().id))]\",\"type\":\"string\"},\"testName\":{\"type\":\"string\"}},\"resources\":[{\"apiVersion\":\"2020-10-01\",\"name\":\"storageAccounts\",\"properties\":{\"expressionEvaluationOptions\":{\"scope\":\"inner\"},\"mode\":\"Incremental\",\"parameters\":{\"location\":{\"value\":\"[resourceGroup().location]\"},\"name\":{\"value\":\"bobjacbicepsa\"}},\"template\":{\"$schema\":\"https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#\",\"contentVersion\":\"1.0.0.0\",\"metadata\":{\"_generator\":{\"name\":\"bicep\",\"templateHash\":\"3993440745275151585\",\"version\":\"0.15.31.15270\"}},\"outputs\":{\"id\":{\"type\":\"string\",\"value\":\"[resourceId('Microsoft.Storage/storageAccounts', parameters('name'))]\"},\"name\":{\"type\":\"string\",\"value\":\"[parameters('name')]\"}},\"parameters\":{\"kind\":{\"allowedValues\":[\"BlobStorage\",\"BlockBlobStorage\",\"FileStorage\",\"Storage\",\"StorageV2\"],\"defaultValue\":\"StorageV2\",\"type\":\"string\"},\"location\":{\"defaultValue\":\"\",\"type\":\"string\"},\"name\":{\"defaultValue\":\"\",\"type\":\"string\"},\"sku_name\":{\"allowedValues\":[\"Premium_LRS\",\"Premium_ZRS\",\"Standard_GRS\",\"Standard_GZRS\",\"Standard_LRS\",\"Standard_RAGRS\",\"Standard_RAGZRS\",\"Standard_ZRS\"],\"defaultValue\":\"Standard_LRS\",\"type\":\"string\"}},\"resources\":[{\"apiVersion\":\"2021-08-01\",\"kind\":\"[parameters('kind')]\",\"location\":\"[parameters('location')]\",\"name\":\"[parameters('name')]\",\"properties\":{\"allowBlobPublicAccess\":false,\"minimumTlsVersion\":\"TLS1_2\",\"networkAcls\":{\"defaultAction\":\"Deny\"},\"publicNetworkAccess\":\"Disabled\"},\"sku\":{\"name\":\"[parameters('sku_name')]\"},\"type\":\"Microsoft.Storage/storageAccounts\"},{\"apiVersion\":\"2021-08-01\",\"dependsOn\":[\"[resourceId('Microsoft.Storage/storageAccounts', parameters('name'))]\"],\"name\":\"[format('{0}/{1}', parameters('name'), 'default')]\",\"properties\":{\"changeFeed\":{\"enabled\":true},\"containerDeleteRetentionPolicy\":{\"days\":181,\"enabled\":true},\"deleteRetentionPolicy\":{\"days\":181,\"enabled\":true},\"isVersioningEnabled\":true,\"restorePolicy\":{\"days\":180,\"enabled\":true}},\"type\":\"Microsoft.Storage/storageAccounts/blobServices\"}]}},\"tags\":{\"modm.events\":\"true\",\"modm.id\":\"31e9f9a0-9fd2-4294-a0a3-0101246d9700\",\"modm.name\":\"storageAccounts\",\"modm.retry\":\"3\",\"modm.stage.id\":\"31e9f9a0-9fd2-4294-a0a3-0101246d9700\"},\"type\":\"Microsoft.Resources/deployments\"}]}"
	err := json.Unmarshal([]byte(deploymentJsonString), &template)
    if err != nil {
		assert.NoError(t, err) 
    }

	parameters := template["parameters"].(map[string]any)
	assert.NotNil(t, parameters)
	assert.Equal(t, len(parameters), 2)
}