# Commercial Marketplace Offer Deployment Manager (MODM)

MODM Development CLI library


development use example:

```
# packaged using 1st party VM offer
modm package build \
    --name "simple terraform app" \
    --description "Simple Terraform application template that deploys a storage account" \
    --version v2.0.0
    --template-file build/managedapp/terraform/simple/templates/main.tf \
    --create-ui-definition build/managedapp/terraform/simple/createUiDefinition.json \
    --out-dir ./bin

# packaged using vm reference
modm package build \
    --name "simple terraform app" \
    --description "Simple Terraform application template that deploys a storage account" \
    --version v2.0.0 \
    --vmi-reference true \
    --template-file build/managedapp/terraform/simple/templates/main.tf \
    --create-ui-definition build/managedapp/terraform/simple/createUiDefinition.json \
    --out-dir ./bin


```