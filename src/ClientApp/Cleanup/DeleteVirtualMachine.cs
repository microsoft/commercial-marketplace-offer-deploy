using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Compute;
using Azure.ResourceManager.Compute.Models;
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

            //bool success = await DisassociateNics(vm);
            //if (!success)
            //{
            //    return new DeleteResourceResult { Succeeded = false };
            //}    

            var operation = await vm.DeleteAsync(WaitUntil.Started);
            var completion = await operation.WaitForCompletionResponseAsync();

            Logger.LogInformation($"The Delete Virtual Machine handler operation completed - status:{completion.IsError}");

            if (completion.IsError)
            {
                Logger.LogError("Deletion of resource {id} failed with status {status}. Reason: {reason}",
                    vm.Id, completion.Status, completion.ReasonPhrase);

                return new DeleteResourceResult { Succeeded = false };
            }

            return new DeleteResourceResult { Succeeded = true };
        }

        private async Task<bool> DisassociateNics(VirtualMachineResource virtualMachine)
        {
            if (virtualMachine.Data.NetworkProfile.NetworkInterfaces.Count > 0)
            {
                var updateOptions = new VirtualMachinePatch()
                {
                    NetworkProfile = new VirtualMachineNetworkProfile()
                };

                foreach (var nic in virtualMachine.Data.NetworkProfile.NetworkInterfaces)
                {
                    updateOptions.NetworkProfile.NetworkInterfaces.Remove(nic);
                }

                var updateOperation = await virtualMachine.UpdateAsync(WaitUntil.Completed, updateOptions);
                var completion = await updateOperation.WaitForCompletionResponseAsync();
                return (!completion.IsError);
            }

            return true;
        }
    }
}