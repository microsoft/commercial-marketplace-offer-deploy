using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Compute;
using Azure.ResourceManager.Resources;

namespace ClientApp.Cleanup
{
    public class DeleteVirtualMachine : IDeleteResourceRequest
    {
        public ResourceIdentifier ResourceId { get; }

        public ResourceGroupResource ResourceGroup { get; }

        public DeleteVirtualMachine(ResourceGroupResource resourceGroup, ResourceIdentifier identifier)
        {
            ResourceGroup = resourceGroup;
            ResourceId = identifier;
        }
    }

    [RetryPolicy]
    public class DeleteVirtualMachineHandler : DeleteResourceHandler<DeleteVirtualMachine>
    {
        public DeleteVirtualMachineHandler(ILoggerFactory loggerFactory, ArmClient client) : base(loggerFactory, client)
        {
        }

        protected override async Task<DeleteResourceResult> DeleteAsync(ResourceGroupResource resourceGroup, ResourceIdentifier id)
        {
            Response<VirtualMachineResource> response = await resourceGroup.GetVirtualMachineAsync(id.Name);
            var vm = response.Value;

            // TODO: disassociate nic

            var operation = await vm.DeleteAsync(WaitUntil.Started);
            var completion = await operation.WaitForCompletionResponseAsync();

            if (completion.IsError)
            {
                Logger.LogError("Deletion of resource {id} failed with status {status}. Reason: {reason}",
                    vm.Id, completion.Status, completion.ReasonPhrase);

                return new DeleteResourceResult { Succeeded = false };
            }

            return new DeleteResourceResult { Succeeded = true };
        }
    }
}