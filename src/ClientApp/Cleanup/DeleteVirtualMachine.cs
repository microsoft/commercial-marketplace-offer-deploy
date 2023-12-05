using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;

namespace ClientApp.Cleanup
{
    public class DeleteVirtualMachine : IDeleteResourceRequest
    {
        public ResourceIdentifier ResourceId { get; }

        public DeleteVirtualMachine(ResourceIdentifier identifier)
        {
            ResourceId = identifier;
        }
    }

    [RetryPolicy(RetryCount = 3, SleepDuration = 1000)]
    public class DeleteVirtualMachineHandler : DeleteResourceHandler<DeleteVirtualMachine>
    {
        public DeleteVirtualMachineHandler(ILoggerFactory loggerFactory, ArmClient client) : base(loggerFactory, client)
        {
        }

        public override async Task<DeleteResourceResult> DeleteAsync(GenericResource resource)
        {
            var operation = await resource.DeleteAsync(WaitUntil.Started);
            var response = await operation.WaitForCompletionResponseAsync();

            if (response.IsError)
            {
                Logger.LogError("Deletion of resource {id} failed with status {status}. Reason: {reason}",
                    resource.Id, response.Status, response.ReasonPhrase);

                return new DeleteResourceResult { Succeeeded = false };
            }

            return new DeleteResourceResult { Succeeeded = true };
        }
    }
}

