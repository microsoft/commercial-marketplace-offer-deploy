
using Azure;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using Microsoft.AspNetCore.Mvc;
using Modm;
using Modm.Engine;

namespace WebHost.Deployments
{
    [Route("api/[controller]")]
    [ApiController]
    public class ResourcesController : ControllerBase
    {
        private readonly ArmClient client;

        public ResourcesController(ArmClient client)
        {
            this.client = client;
        }

        [HttpGet]
        [Route("{resourceGroupName}")]
        public async Task<List<string>> GetAsync([FromRoute] string resourceGroupName)
        {
            var subscription = await client.GetDefaultSubscriptionAsync();
            var response = await subscription.GetResourceGroupAsync(resourceGroupName);
            var resourceGroup = response.Value;

            var resources = new List<string>();

            await foreach (var resource in resourceGroup.GetGenericResourcesAsync())
            {
                resources.Add(resource.Id.ToString());
            }

            return resources;
        }
    }
}
