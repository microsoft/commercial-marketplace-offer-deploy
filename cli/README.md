# Commercial Marketplace Offer Deployment Manager (MODM)

MODM Development CLI library


development use example:

```
modm package build \
    --name "simple terraform app" \
    --description "Simple Terraform application template that deploys a storage account" \
    --vmi-reference-id vmi-reference-id \
    --template-file build/managedapp/terraform/simple/templates/main.tf \
    --create-ui-definition build/managedapp/terraform/simple/createUiDefinition.json \
    --out-dir ./bin

```