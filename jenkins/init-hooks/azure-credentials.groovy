#!groovy

// This script will add the Azure Credentials that can be used within Jobs
//      it will select either a ServicePrincipal or Managed Identity credential 
//      based on the presence of environment variables present

import com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl
import org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl
import org.jenkinsci.plugins.plaincredentials.impl.FileCredentialsImpl
import com.cloudbees.plugins.credentials.domains.Domain
import com.cloudbees.plugins.credentials.CredentialsScope
import jenkins.model.Jenkins
import com.microsoft.azure.util.AzureImdsCredentials
import com.microsoft.azure.util.AzureCredentials

final SYSTEM_CREDENTIALS_PROVIDER = 'com.cloudbees.plugins.credentials.SystemCredentialsProvider'

// use this ID to reference and configure credential usage in any jenkins job to acquire credentials
final CREDENTIALS_ID = '59aaa22a-e04e-4909-9cc8-2ac406e002d0'

def shouldCredentialsBeServicePrincipal(clientId, clientSecret) {
        return (clientId != null && clientId != '') && (clientSecret != null && clientSecret != '')
}

def instance = Jenkins.get()
def domain = Domain.global()
def store = instance.getExtensionList(SYSTEM_CREDENTIALS_PROVIDER)[0].getStore()

final subscriptionId = System.getenv('AZURE_SUBSCRIPTION_ID')
final clientId = System.getenv('AZURE_CLIENT_ID')
final clientSecret = System.getenv('AZURE_CLIENT_SECRET')

// if there's a client ID and a client secret present, then it should be a service principal
// otherwise, default to managed identity
if (shouldCredentialsBeServicePrincipal(clientId, clientSecret)) {
        final description = 'Azure Service Principal Credentials'
        def servicePrincipalCredentials = new AzureCredentials(
                CredentialsScope.GLOBAL, 
                CREDENTIALS_ID, 
                description,
                subscriptionId,
                clientId,
                clientSecret)
        store.addCredentials(domain, servicePrincipalCredentials)

} else { // managed identity
        // the env name is the target cloud type / scope for azure, e.g. commercial, Gov, etc. it has nothing to do with env variables
        final azureEnvironmentName = 'Azure'
        def managedIdentityCredentials = new AzureImdsCredentials(
                CredentialsScope.GLOBAL, 
                CREDENTIALS_ID, 
                azureEnvironmentName)

        managedIdentityCredentials.subscriptionId = subscriptionId
        managedIdentityCredentials.clientId = clientId
        store.addCredentials(domain, managedIdentityCredentials)
}