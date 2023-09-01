#!/bin/bash

echo "Hello crom deploy.sh"

cd terraformSol

# Initialize Terraform (required before first run)
terraform init -backend=false

# Apply the parent Terraform template
terraform apply -auto-approve
