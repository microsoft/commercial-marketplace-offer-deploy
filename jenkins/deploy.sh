#!/bin/bash

echo "Hello crom deploy.sh"

cd $JENKINS_HOME/solutions/terraform/content

# Initialize Terraform (required before first run)
terraform init -backend=false

# Set Azure credentials from Jenkins bindings
export ARM_CLIENT_ID=$AZURE_CLIENT_ID
export ARM_CLIENT_SECRET=$AZURE_CLIENT_SECRET
export ARM_SUBSCRIPTION_ID=$AZURE_SUBSCRIPTION_ID
export ARM_TENANT_ID=$AZURE_TENANT_ID

# Apply the parent Terraform template
terraform apply -auto-approve
