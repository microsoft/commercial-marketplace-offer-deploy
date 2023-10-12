#!/bin/bash

# ===========================================================================================
#   DESCRIPTION:
#     Provisioning shell script, executed during Packer build
#
#   Execution context:
#     This script is executed on the VM by Packer, not locally
# ===========================================================================================


echo set debconf to Noninteractive
echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections


# prep with prerequisites
sudo apt-get update -y
sudo apt-get upgrade -y

# Docker engine
sudo apt-get install ca-certificates curl gnupg -y

# Add Docker GPG key and repository
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Update package lists after adding the repository
sudo apt-get update -y

# Install Docker packages
sudo apt-get install docker-ce docker-ce-cli containerd.io -y

# Optional: Install additional Docker-related packages
sudo apt-get install docker-buildx-plugin docker-compose -y

# Install .NET 7
sudo apt-get install -y dotnet-sdk-7.0

echo "Installing Azure CLI"

# Install Azure Functions Core needed for the Azure Function App
# This will be used to publish the dashboard redirect
curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.gpg
sudo mv microsoft.gpg /etc/apt/trusted.gpg.d/microsoft.gpg

sudo sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/repos/microsoft-ubuntu-$(lsb_release -cs)-prod $(lsb_release -cs) main" > /etc/apt/sources.list.d/dotnetdev.list'

sudo apt-get update
sudo apt-get install azure-functions-core-tools-4


# clone the MODM source into source
sudo git clone --depth=1 --branch $MODM_REPO_BRANCH https://github.com/microsoft/commercial-marketplace-offer-deploy.git $MODM_HOME/source

# Initial image setup so we can get cached image layers to speed up builds for the final vmi
cd $MODM_HOME/source

echo "Building docker images."

sudo docker build ./src -t modm -f ./build/container/Dockerfile.modm  
sudo docker build . -t jenkins -f ./build/container/Dockerfile.jenkins
