#! /bin/bash

# DESCRIPTION
# Creates an application registration on the target Azure Active Directory and configures it.
#
# Purposes:
#   Use this to setup an app registration for local development
#   1. login using az login
#   2. Execute this script to create an app registration


app_registration_name="ApiServer/DevTest"

client_id=$(az ad app create --display-name $app_registration_name --query appId --output tsv)
object_id=$(az ad app show --id $client_id --query id --output tsv)

# remove Azure Active Directory Graph permission and Azure Active Directory Graph delegated permission
az ad app permission delete --id $client_id --api 00000002-0000-0000-c000-000000000000
az ad app permission delete --id $client_id --api 00000002-0000-0000-c000-000000000000 --api-permissions 311a71cc-e848-46a1-bdf8-97ff7156d8e6

# add the API identifier URI
az ad app update --id $client_id --identifier-uris "api://marketplaceofferdeploy"

# add scope
scope_id=$(uuidgen | awk '{print tolower($0)}')
oauth2_permissions=$(cat <<EOF
{
    "oauth2Permissions": [
		{
			"adminConsentDescription": "Default Access",
			"adminConsentDisplayName": "Default Access",
			"id": "$scope_id",
			"isEnabled": true,
			"lang": null,
			"origin": "Application",
			"type": "User",
			"userConsentDescription": "Default Access",
			"userConsentDisplayName": "Default Access",
			"value": "api_access"
		}
	]
}
EOF
)
az ad app update --id $client_id --required-resource-accesses $oauth2_permissions

# add required resource access (the scope)
required_resource_access=$(cat <<EOF
{
 "requiredResourceAccess": [
     {
         "resourceAppId": "da72dd07-7708-4ccf-a567-fabc36c0edbf",
         "resourceAccess": [
             {
                 "id": "69c3d80f-e337-467e-862f-b3ed1a4bedc4",
                 "type": "Scope"
             }
         ]
     }
 ]
}
EOF
)
az ad app update --id $client_id --required-resource-accesses $required_resource_access

# add azure cli to the authorized clients
azure_cli_client_id=04b07795-8ddb-461a-bbee-02f9e1bf7b46

pre_authorized_applications=$(cat <<EOF
{
    "preAuthorizedApplications": [
		{
			"appId": "$azure_cli_client_id",
			"permissionIds": [
				"69c3d80f-e337-467e-862f-b3ed1a4bedc4"
			]
		}
	]
}
EOF
)
az ad app update --id $client_id --set api="$pre_authorized_applications"
