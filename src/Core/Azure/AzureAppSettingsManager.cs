// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using Microsoft.Azure.Management.Fluent;
using Microsoft.Azure.Management.ResourceManager.Fluent;
using Microsoft.Azure.Management.ResourceManager.Fluent.Authentication;
using System.Collections.Generic;
using System.Threading.Tasks;
using Azure.Core;
using Azure.Identity;
using Microsoft.Rest;

namespace Modm.Azure
{
    public class AzureAppSettingsManager
    {
        private readonly string subscriptionId;
        private readonly DefaultAzureCredential credential;
        private readonly Microsoft.Azure.Management.Fluent.Azure.IAuthenticated authenticatedAzure;

        public AzureAppSettingsManager(string subscriptionId)
        {
            this.subscriptionId = subscriptionId;
            this.credential = new DefaultAzureCredential();

            // Setup Azure Fluent SDK Authentication
            var tokenCredentials = new TokenCredentials(credential.GetToken(new TokenRequestContext(new[] { "https://management.azure.com/.default" })).Token);
            var azureCredentials = new AzureCredentials(tokenCredentials, tokenCredentials, string.Empty, AzureEnvironment.AzureGlobalCloud);
            this.authenticatedAzure = Microsoft.Azure.Management.Fluent.Azure.Authenticate(azureCredentials);
        }

        public async Task UpdateAppSettingsAsync(string resourceGroupName, string appName, Dictionary<string, string> newSettings)
        {
            var azure = authenticatedAzure.WithSubscription(subscriptionId);
            var appService = await azure.AppServices.WebApps.GetByResourceGroupAsync(resourceGroupName, appName);

            foreach (var setting in newSettings)
            {
                appService.Update().WithAppSetting(setting.Key, setting.Value).Apply();
            }
        }
    }
}

