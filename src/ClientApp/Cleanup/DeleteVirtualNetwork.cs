using System;
using Azure;
using Azure.Core;
using Azure.Identity;
using Azure.ResourceManager;
using Azure.ResourceManager.Compute;
using Azure.ResourceManager.Network;
using Azure.ResourceManager.Resources;
using Microsoft.CodeAnalysis;
using Microsoft.Identity.Client.Platforms.Features.DesktopOs.Kerberos;

namespace ClientApp.Cleanup
{
    public class DeleteVirtualNetwork : IDeleteResourceRequest
    {
        public ResourceIdentifier ResourceId { get; }

        public ResourceGroupResource ResourceGroup { get; }

        public DeleteVirtualNetwork(ResourceGroupResource resourceGroup, ResourceIdentifier identifier)
        {
            ResourceGroup = resourceGroup;
            ResourceId = identifier;
        }
    }

    [RetryPolicy]
    public class DeleteVirtualNetworkHandler : DeleteResourceHandler<DeleteVirtualNetwork>
    {
        public DeleteVirtualNetworkHandler(ILoggerFactory loggerFactory, ArmClient client) : base(loggerFactory, client)
        {

        }

        protected override async Task<DeleteResourceResult> DeleteAsync(ResourceGroupResource resourceGroup, ResourceIdentifier id)
        {
            var response = await resourceGroup.GetVirtualNetworkAsync(id.Name);
            var vnet = response.Value;
            var operation = await vnet.DeleteAsync(WaitUntil.Started);
            var completion = await operation.WaitForCompletionResponseAsync();

            if (completion.IsError)
            {
                Logger.LogError("Deletion of resource {id} failed with status {status}. Reason: {reason}",
                vnet.Id, completion.Status, completion.ReasonPhrase);

                return new DeleteResourceResult { Succeeded = false };
            }

            return new DeleteResourceResult { Succeeded = true };
        }
    }
}