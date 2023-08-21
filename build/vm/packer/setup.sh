#!/bin/bash

subscription_id=$(az account show --query id --output tsv)

az group create -n modmImage -l eastus
az ad sp create-for-rbac --name "modmImage" --role Contributor --scopes /subscriptions/$subscription_id --query "{ client_id: appId, client_secret: password, tenant_id: tenant }"
az account show --query "{ subscription_id: id }"

