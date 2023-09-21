using Azure.ResourceManager;

namespace Modm.Deployments
{
    public class DeploymentResourcesClient
	{
        private readonly ArmClient client;

        public DeploymentResourcesClient(ArmClient client)
		{
            this.client = client;
        }

        public async Task<IEnumerable<DeploymentResource>> Get(string resourceGroupName)
        {
            try
            {
                var subscription = await client.GetDefaultSubscriptionAsync();
                var resourceGroup = await subscription.GetResourceGroupAsync(resourceGroupName);
                var resources = await resourceGroup.Value.GetGenericResourcesAsync(
                    expand: "provisioningState,"
                    ).ToListAsync();

                return resources.Select(r => new DeploymentResource
                {
                    Name = r.Data.Name,
                    Type = r.Data.ResourceType.ToString(),
                    State = r.Data.ProvisioningState,
                    Timestamp = r.Data.CreatedOn.GetValueOrDefault(DateTimeOffset.UtcNow)
                });
            }
            catch
            {
                return new List<DeploymentResource>();
            }
        }
	}
}

