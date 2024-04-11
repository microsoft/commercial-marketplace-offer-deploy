using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using Azure.ResourceManager.Resources.Models;
using Microsoft.AspNetCore.Mvc;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace WebHost.WebHost.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class ResourcesController : ControllerBase
    {
        private readonly ArmClient client;
        private readonly ILogger<ResourcesController> logger;

        public ResourcesController(ArmClient client, ILogger<ResourcesController> logger)
        {
            this.client = client;
            this.logger = logger;
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
