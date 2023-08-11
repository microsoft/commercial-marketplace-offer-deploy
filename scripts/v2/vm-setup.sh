#!/bin/bash


# Docker engine
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg -y

sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
sudo apt install git -y

# clone repository against develop branch
git clone --depth=1 --branch develop https://github.com/microsoft/commercial-marketplace-offer-deploy.git
cd commercial-marketplace-offer-deploy

echo SITE_ADDRESS=$SITE_ADDRESS > ./deployments/.env
sudo docker compose -f ./deployments/docker-compose.development.yml -p modm up