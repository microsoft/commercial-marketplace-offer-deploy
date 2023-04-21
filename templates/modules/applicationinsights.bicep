@description('Name of Application Insights resource.')
param name string = 'modm-appinsights'

@description('Type of app you are deploying. This field is for legacy reasons and will not impact the type of App Insights resource you deploy.')
param type string = 'web'

@description('Location for all resources.')
param location string = resourceGroup().location

@description('Source of Azure Resource Manager deployment')
param requestSource string = 'rest'

resource component 'Microsoft.Insights/components@2020-02-02' = {
  name: name
  location: location
  kind: 'other'
  properties: {
    Application_Type: type
    Flow_Type: 'Bluefield'
    Request_Source: requestSource
  }
}
