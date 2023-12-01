using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Modm.Deployments;
using Modm.Diagnostics;
using Modm.Engine;
using ClientApp.Backend;
using Azure.ResourceManager;
using Modm.Azure;
using MediatR;
using ClientApp.Notifications;

namespace Modm.ClientApp.Controllers
{
    [Route("api")]
    [ApiController]
    [Authorize]
    public class ProxyController : ControllerBase
    {
        private readonly ProxyClientFactory clientFactory;
        private IProxyClient client;

        private readonly IAzureResourceManager resourceManager;
        private readonly IMediator mediator;
        private readonly ILogger<ProxyController> logger;

        /// <summary>
        /// The instance of the proxy client based on the incoming http request
        /// </summary>
        private IProxyClient Client
        {
            get { return client ??= clientFactory.Create(HttpContext.Request); }
        }

        public ProxyController(ProxyClientFactory clientFactory, IMediator mediator, IAzureResourceManager resourceManager, ILogger<ProxyController> logger)
        {
            this.clientFactory = clientFactory;
            this.mediator = mediator;
            this.resourceManager = resourceManager;
            this.logger = logger;
        }

        [HttpPost]
        [Route("resources/{resourceGroupName}/deletemodmresources")]
        public async Task<IActionResult> DeleteResourcesWithTagAsync([FromRoute] string resourceGroupName)
        {
            var deleteInitiated = new DeleteInitiated(resourceGroupName);
            await this.mediator.Send(deleteInitiated);

            return Ok("Successfully submitted a delete");
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
