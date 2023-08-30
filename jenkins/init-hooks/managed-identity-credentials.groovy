#!groovy

import com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl
import org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl
import org.jenkinsci.plugins.plaincredentials.impl.FileCredentialsImpl
import com.cloudbees.plugins.credentials.domains.Domain
import com.cloudbees.plugins.credentials.CredentialsScope
import jenkins.model.Jenkins
import com.microsoft.azure.util.AzureImdsCredentials


final SYSTEM_CREDENTIALS_PROVIDER = 'com.cloudbees.plugins.credentials.SystemCredentialsProvider'

// use this ID to reference and configure credential usage in any jenkins job
// to acquire credentials
final MANAGED_IDENTITY_CREDENTIALS_ID = '59aaa22a-e04e-4909-9cc8-2ac406e002d0'
final MANAGED_IDENTITY_CREDENTIALS_ENV_NAME = 'Azure'

def instance = Jenkins.get()
def store = instance.getExtensionList(SYSTEM_CREDENTIALS_PROVIDER)[0].getStore()


def managedIdentityCredentials = new AzureImdsCredentials(
        CredentialsScope.GLOBAL, 
        MANAGED_IDENTITY_CREDENTIALS_ID, 
        MANAGED_IDENTITY_CREDENTIALS_ENV_NAME
)

managedIdentityCredentials.subscriptionId = System.getenv('AZURE_SUBSCRIPTION_ID')
managedIdentityCredentials.clientId = System.getenv('AZURE_CLIENT_ID')

// add to the global store
def domain = Domain.global()
store.addCredentials(domain, managedIdentityCredentials)