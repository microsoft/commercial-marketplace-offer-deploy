#!/bin/bash

# prep with prerequisites
sudo apt-get update
sudo apt-get upgrade -y

# Docker engine
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg -y

# Add Docker GPG key and repository
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Update package lists after adding the repository
sudo apt-get update

# Install Docker packages
sudo apt-get install docker-ce docker-ce-cli containerd.io -y

# Optional: Install additional Docker-related packages
sudo apt-get install docker-buildx-plugin docker-compose -y

# Install .NET 7
sudo apt-get install -y dotnet-sdk-7.0