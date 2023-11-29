using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Modm.Deployments;
using Modm.Diagnostics;
using Modm.Engine;
using ClientApp.Backend;
using Azure.ResourceManager;
using Modm.Azure;

namespace Modm.ClientApp.Controllers
{
    [Route("api")]
    [ApiController]
    [Authorize]
    public class ProxyController : ControllerBase
    {
        private readonly ProxyClientFactory clientFactory;
        private IProxyClient client;

        private readonly ArmClient armClient;
        private readonly ILogger<ProxyController> logger;

        /// <summary>
        /// The instance of the proxy client based on the incoming http request
        /// </summary>
        private IProxyClient Client
        {
            get { return client ??= clientFactory.Create(HttpContext.Request); }
        }

        public ProxyController(ProxyClientFactory clientFactory, ArmClient armClient, ILogger<ProxyController> logger)
        {
            this.clientFactory = clientFactory;
            this.armClient = armClient;
            this.logger = logger;
        }

        [HttpPost]
        [Route("resources/{resourceGroupName}/deletemodmresources")]
        public async Task<IActionResult> DeleteResourcesWithTagAsync([FromRoute] string resourceGroupName)
        {
            var armCleanup = new AzureDeploymentCleanup(armClient);
            bool deleted = await armCleanup.DeleteResourcePostDeployment(resourceGroupName);

            if (!deleted)
            {
                return BadRequest("Some resources could not be deleted after multiple attempts.");
            }

            return Ok("Resources with modm tag applied have been deleted. App Service Being deleted without waiting.");
        }


        [HttpGet("deployments")]
        public async Task<IActionResult> GetDeployments()
        {
            return await Client.GetAsync<GetDeploymentResponse>(Routes.GetDeployments);
        }

        [HttpGet("diagnostics")]
        public async Task<IActionResult> GetDiagnostics()
        {
            return await Client.GetAsync<GetDiagnosticsResponse>(Routes.GetDiagnostics);
        }

        [HttpGet("status")]
        public async Task<IActionResult> GetStatus()
        {
            return await Client.GetAsync<EngineInfo>(Routes.GetStatus);
        }
    }
}
