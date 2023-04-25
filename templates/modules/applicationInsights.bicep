@description('MODM app version')
param appVersion string

@description('Location.')
param location string = resourceGroup().location

var versionSuffix = replace(appVersion, '.', '')
var name = 'modm-appinsights-${versionSuffix}'

resource appInsights 'Microsoft.Insights/components@2020-02-02' = {
  name: name
  location: location
  kind: 'other'
  properties: {
    Application_Type: 'web'
    Flow_Type: 'Bluefield'
    Request_Source: 'rest'
  }
}

output appInsightsName string = appInsights.name
output appInsightsInstrumentationKey string = appInsights.properties.InstrumentationKey
