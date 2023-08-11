#!/bin/bash

# creates a virtual machine in azure

export RESOURCE_GROUP_NAME=modm-dev
export VNET_NAME=modm-vnet-dev
export LOCATION=eastus
export VM_NAME=modm-vm02-dev
export VM_IMAGE=debian
export ADMIN_USERNAME=azureuser
export VM_DNS_NAME=modmvm02dev
export VM_FQDN=$VM_DNS_NAME.eastus.cloudapp.azure.com

az vm create \
  --resource-group $RESOURCE_GROUP_NAME \
  --name $VM_NAME \
  --image $VM_IMAGE \
  --assign-identity [system] \
  --vnet-name $VNET_NAME \ 
  --public-ip-address-dns-name $VM_DNS_NAME \
  --admin-username $ADMIN_USERNAME \
  --generate-ssh-keys \
  --public-ip-sku Standard \
  --size Standard_D2s_v3


export IP_ADDRESS=$(az vm show --show-details --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --query publicIps --output tsv)
az vm open-port --port 443 --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME


az vm run-command invoke \
   --resource-group $RESOURCE_GROUP_NAME \
   --name $VM_NAME \
   --command-id RunShellScript \
   --scripts @vm-setup.sh

az vm run-command invoke \
   --resource-group $RESOURCE_GROUP_NAME \
   --name $VM_NAME \
   --command-id RunShellScript \
   --scripts "docker compose -f $compose_file -p modm up"


echo "Access at https://$VM_FQDN"