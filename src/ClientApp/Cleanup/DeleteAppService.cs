// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.AppService;
using Azure.ResourceManager.Resources;

namespace ClientApp.Cleanup
{
    public class DeleteAppService : IDeleteResourceRequest
    {
        public ResourceIdentifier ResourceId { get; }

        public ResourceGroupResource ResourceGroup { get; }

        public DeleteAppService(ResourceGroupResource resourceGroup, ResourceIdentifier identifier)
        {
            ResourceGroup = resourceGroup;
            ResourceId = identifier;
        }
    }

    [RetryPolicy]
    public class DeleteAppServiceHandler : DeleteResourceHandler<DeleteAppService>
    {
        public DeleteAppServiceHandler(ILoggerFactory loggerFactory, ArmClient client) : base(loggerFactory, client)
        {

        }

        protected override async Task<DeleteResourceResult> DeleteAsync(ResourceGroupResource resourceGroup, ResourceIdentifier id)
        {
            var appServiceResponse = await resourceGroup.GetWebSiteAsync(id.Name);
            var appService = appServiceResponse.Value;

            var delete = await appService.DeleteAsync(WaitUntil.Started);
            var completion = await delete.WaitForCompletionResponseAsync();

            if (completion.IsError)
            {
                Logger.LogError("Deletion of resource {id} failed with status {status}. Reason: {reason}",
                appService.Id, completion.Status, completion.ReasonPhrase);

                return new DeleteResourceResult { Succeeded = false };
            }

            return new DeleteResourceResult { Succeeded = true };
        }
    }
}

