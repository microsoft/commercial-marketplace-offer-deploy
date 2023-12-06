using System;
using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Storage;
using Azure.ResourceManager.Resources;


namespace ClientApp.Cleanup
{
    public class DeleteStorageAccount : IDeleteResourceRequest
    {
        public ResourceIdentifier ResourceId { get; }

        public ResourceGroupResource ResourceGroup { get; }

        public DeleteStorageAccount(ResourceGroupResource resourceGroup, ResourceIdentifier identifier)
        {
            ResourceGroup = resourceGroup;
            ResourceId = identifier;
        }
    }

    [RetryPolicy]
    public class DeleteStorageAccountHandler : DeleteResourceHandler<DeleteStorageAccount>
    {
        public DeleteStorageAccountHandler(ILoggerFactory loggerFactory, ArmClient client) : base(loggerFactory, client)
        {

        }

        protected async override Task<DeleteResourceResult> DeleteAsync(ResourceGroupResource resourceGroup, ResourceIdentifier id)
        {
            var storageAccountResponse = await resourceGroup.GetStorageAccountAsync(id.Name);
            var storageAccount = storageAccountResponse.Value;

            var delete = await storageAccount.DeleteAsync(WaitUntil.Started);
            var completion = await delete.WaitForCompletionResponseAsync();

            if (completion.IsError)
            {
                Logger.LogError("Deletion of resource {id} failed with status {status}. Reason: {reason}",
                storageAccount.Id, completion.Status, completion.ReasonPhrase);

                return new DeleteResourceResult { Succeeded = false };
            }

            return new DeleteResourceResult { Succeeded = true };
        }
    }
}