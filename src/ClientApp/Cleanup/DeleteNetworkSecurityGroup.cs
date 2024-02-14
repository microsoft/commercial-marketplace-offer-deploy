// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Network;
using Azure.ResourceManager.Resources;

namespace ClientApp.Cleanup
{
    public class DeleteNetworkSecurityGroup : IDeleteResourceRequest
    {
        public ResourceIdentifier ResourceId { get; }

        public ResourceGroupResource ResourceGroup { get; }

        public DeleteNetworkSecurityGroup(ResourceGroupResource resourceGroup, ResourceIdentifier identifier)
        {
            ResourceGroup = resourceGroup;
            ResourceId = identifier;
        }
    }

    [RetryPolicy]
    public class DeleteNetworkSecurityGroupHandler : DeleteResourceHandler<DeleteNetworkSecurityGroup>
    {
        public DeleteNetworkSecurityGroupHandler(ILoggerFactory loggerFactory, ArmClient client) : base(loggerFactory, client)
        {

        }

        protected override async Task<DeleteResourceResult> DeleteAsync(ResourceGroupResource resourceGroup, ResourceIdentifier id)
        {
            var response = await resourceGroup.GetNetworkSecurityGroupAsync(id.Name);
            var nsg = response.Value;
            var operation = await nsg.DeleteAsync(WaitUntil.Started);
            var completion = await operation.WaitForCompletionResponseAsync();

            if (completion.IsError)
            {
                Logger.LogError("Deletion of resource {id} failed with status {status}. Reason: {reason}",
                nsg.Id, completion.Status, completion.ReasonPhrase);

                return new DeleteResourceResult { Succeeded = false };
            }

            return new DeleteResourceResult { Succeeded = true };
        }
    }
}
