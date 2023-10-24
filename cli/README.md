# Commercial Marketplace Offer Deployment Manager (MODM)

MODM Development CLI library


development use example:

```
# packaged using 1st party VM offer
modm package build \
    --name "simple terraform app" \
    --description "Simple Terraform application template that deploys a storage account" \
    --version v2.0.0 \
    --template-file build/managedapp/terraform/simple/templates/main.tf \
    --create-ui-definition build/managedapp/terraform/simple/createUiDefinition.json \
    --out-dir ./bin

# packaged using vmi reference
modm package build \
    --name "simple terraform app" \
    --description "Simple Terraform application template that deploys a storage account" \
    --version v2.0.0 \
    --vmi-reference true \
    --template-file build/managedapp/terraform/simple/templates/main.tf \
    --create-ui-definition build/managedapp/terraform/simple/createUiDefinition.json \
    --out-dir ./bin

# packaged using vmi reference id that will be used directly
modm package build \
    --name "simple terraform app" \
    --description "Simple Terraform application template that deploys a storage account" \
    --version v2.0.0 \
    --vmi-reference-id /subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/modm-dev-vmi/providers/Microsoft.Compute/galleries/modm.dev.sig/images/modm/versions/0.1.96 \
    --template-file build/managedapp/terraform/simple/templates/main.tf \
    --create-ui-definition build/managedapp/terraform/simple/createUiDefinition.json \
    --out-dir ./bin
```