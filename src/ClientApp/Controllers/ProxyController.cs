using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Modm.Deployments;
using Modm.Diagnostics;
using Modm.Engine;
using ClientApp.Backend;

namespace Modm.ClientApp.Controllers
{
    [Route("api")]
    [ApiController]
    [Authorize]
    public class ProxyController : ControllerBase
    {
        private readonly ProxyClientFactory clientFactory;
        private IProxyClient client;

        /// <summary>
        /// The instance of the proxy client based on the incoming http request
        /// </summary>
        private IProxyClient Client
        {
            get { return client ??= clientFactory.Create(HttpContext.Request); }
        }

        public ProxyController(ProxyClientFactory clientFactory)
        {
            this.clientFactory = clientFactory;
        }

        [HttpPost]
        [Route("resources/{resourceGroupName}/deletemodmresources")]
        public async Task<IActionResult> DeleteResourcesWithTagAsync([FromRoute] string resourceGroupName)
        {
            var relativeUri = string.Format(Routes.DeleteInstallerFormat, resourceGroupName);
            return await Client.PostAsync(relativeUri);
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
