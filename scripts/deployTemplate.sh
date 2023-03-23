RESOURCE_GROUP_NAME="MODMTest"
DEPLOYMENT_NAME="bobjac1"

az deployment group create --resource-group $RESOURCE_GROUP_NAME --name $DEPLOYMENT_NAME --parameters @parameters.json --template-file mainTemplate.json

#### Command to deploy the deployment-orchestrator when not using the ARM deployment (assumes latest app is in local orchestrator-web.zip)
# az webapp deploy -g neudesic-test-deployment-rg -n GEN-UNIQUE-121Portal --src-path orchestrator-web.zip