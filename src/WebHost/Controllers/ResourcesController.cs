using Azure.ResourceManager;
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

        [HttpPost]
        [Route("{resourceGroupName}/deletemodmresources")]
        public async Task<IActionResult> DeleteResourcesWithTagAsync([FromRoute] string resourceGroupName)
        {
            var subscription = await client.GetDefaultSubscriptionAsync();
            var response = await subscription.GetResourceGroupAsync(resourceGroupName);
            var resourceGroup = response.Value;

            await foreach (var resource in resourceGroup.GetGenericResourcesAsync())
            {
                if (resource.Data.Tags != null && resource.Data.Tags.TryGetValue("modm", out var tagValue) && tagValue == "true")
                {
                    await resource.DeleteAsync(Azure.WaitUntil.Started);
                }
            }

            return Ok("Resources with tag modm=true have been deleted.");
        }
    }
}
