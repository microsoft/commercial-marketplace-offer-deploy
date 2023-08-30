#!/bin/bash

## This script is the script executed in the command portion of the Jenkins build step

az login --service-principal -u $AZURE_CLIENT_ID -p $AZURE_CLIENT_SECRET --tenant $AZURE_TENANT_ID
az account set -s $AZURE_SUBSCRIPTION_ID

az account show

echo ""
echo "setting up terraform"

# this path is by convention. MODM will place all artifacts from the content.zip from the artifacts URI
# into a host path that must be mounted to /var/jenkins_home/artifacts
cd $JENKINS_HOME/modm/artifacts

echo "Terraform files"
ls -1

export ARM_CLIENT_ID=$AZURE_CLIENT_ID
export ARM_CLIENT_SECRET=$AZURE_CLIENT_SECRET
export ARM_SUBSCRIPTION_ID=$AZURE_SUBSCRIPTION_ID
export ARM_TENANT_ID=$AZURE_TENANT_ID

terraform init
terraform plan -var-file="../parameters.tfvars" -out main.tfplan
terraform apply -auto-approve main.tfplan