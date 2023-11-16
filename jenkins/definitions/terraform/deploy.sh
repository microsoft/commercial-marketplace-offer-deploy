#!/bin/bash

cd $MODM_HOME/installer

if [ -z "$AZURE_CLIENT_SECRET" ]; then
  az login --identity
  export ARM_USE_MSI=true
  export ARM_CLIENT_ID=$AZURE_CLIENT_ID
  export ARM_SUBSCRIPTION_ID=$AZURE_SUBSCRIPTION_ID
  export ARM_TENANT_ID=$AZURE_TENANT_ID
else
  # Set Azure credentials from Jenkins bindings
  export ARM_CLIENT_ID=$AZURE_CLIENT_ID
  export ARM_CLIENT_SECRET=$AZURE_CLIENT_SECRET
  export ARM_SUBSCRIPTION_ID=$AZURE_SUBSCRIPTION_ID
  export ARM_TENANT_ID=$AZURE_TENANT_ID
fi

# special case where echoing this will strip the Jenkins preamble
# from the job output
echo "-----------------"

# Initialize Terraform (required before first run)
terraform init -backend=false

# Apply the parent Terraform template
terraform apply -auto-approve
