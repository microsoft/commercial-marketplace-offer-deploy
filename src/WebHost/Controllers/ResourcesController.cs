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

        [HttpPost]
        [Route("{resourceGroupName}/deletemodmresources")]
        public async Task<IActionResult> DeleteResourcesWithTagAsync([FromRoute] string resourceGroupName)
        {
            var resourcesToDelete = await GetResourcesToDelete(resourceGroupName);

            int maxAttempts = resourcesToDelete.Count * 5; 
            int attempt = 0;

            while (resourcesToDelete.Count > 0 && attempt < maxAttempts)
            {
                var resource = resourcesToDelete[0]; 

                if (await TryDeleteResource(resource))
                {
                    resourcesToDelete.RemoveAt(0);
                }
                else
                {
                    // If deletion fails, move the resource to the end of the list
                    resourcesToDelete.RemoveAt(0);
                    resourcesToDelete.Add(resource);
                }

                attempt++;
            }

            if (resourcesToDelete.Count == 0)
            {
                return Ok("Resources with tag modm=true have been deleted.");
            }
            else
            {
                return BadRequest("Some resources could not be deleted after multiple attempts.");
            }
        }

        private async Task<List<GenericResource>> GetResourcesToDelete(string resourceGroupName)
        {
            var subscription = await client.GetDefaultSubscriptionAsync();
            var response = await subscription.GetResourceGroupAsync(resourceGroupName);
            var resourceGroup = response.Value;

            var resourcesToDelete = new List<GenericResource>();

            await foreach (var resource in resourceGroup.GetGenericResourcesAsync())
            {
                if (resource.Data.Tags != null && resource.Data.Tags.TryGetValue("modm", out var tagValue) && tagValue == "true")
                {
                    resourcesToDelete.Add(resource);
                }
            }

            return resourcesToDelete;
        }

        private async Task<bool> TryDeleteResource(GenericResource resource)
        {
            try
            {
                await resource.DeleteAsync(Azure.WaitUntil.Started);
                return true;
            }
            catch
            {
                return false; // Return false if deletion fails
            }
        }
    }
}
