// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.AppConfiguration;
using Azure.ResourceManager.Resources;

namespace ClientApp.Cleanup
{
    public class DeleteAppConfiguration : IDeleteResourceRequest
    {
        public ResourceIdentifier ResourceId { get; }

        public ResourceGroupResource ResourceGroup { get; }

        public DeleteAppConfiguration(ResourceGroupResource resourceGroup, ResourceIdentifier identifier)
        {
            ResourceGroup = resourceGroup;
            ResourceId = identifier;
        }
    }

    [RetryPolicy]
    public class DeleteAppConfigurationHandler : DeleteResourceHandler<DeleteAppConfiguration>
    {
        public DeleteAppConfigurationHandler(ILoggerFactory loggerFactory, ArmClient client) : base(loggerFactory, client)
        {
        }

        protected override async Task<DeleteResourceResult> DeleteAsync(ResourceGroupResource resourceGroup, ResourceIdentifier id)
        {
            var appConfigResponse = await resourceGroup.GetAppConfigurationStoreAsync(id.Name);
            var appConfig = appConfigResponse.Value;

            var delete = await appConfig.DeleteAsync(WaitUntil.Started);
            var completion = await delete.WaitForCompletionResponseAsync();

            if (completion.IsError)
            {
                Logger.LogError("Deletion of resource {id} failed with status {status}. Reason: {reason}",
                appConfig.Id, completion.Status, completion.ReasonPhrase);

                return new DeleteResourceResult { Succeeded = false };
            }

            return new DeleteResourceResult { Succeeded = true };
        }
    }
}
