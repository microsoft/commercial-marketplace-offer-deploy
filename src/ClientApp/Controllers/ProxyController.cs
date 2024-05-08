using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Modm.Deployments;
using Modm.Diagnostics;
using Modm.Engine;
using ClientApp.Backend;
using Azure.ResourceManager;
using Modm.Azure;
using MediatR;
using ClientApp.Commands;
using System.Text;
using System.Text.Json;

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

        [HttpPost("deployments/{deploymentId}/redeploy")]
        public async Task<IActionResult> PostRedeploy(int deploymentId, [FromBody] Dictionary<string, object> parameters)
        {
            var request = new StartRedeploymentRequest
            {
                DeploymentId = deploymentId,
                Parameters = parameters
            };

            // Use the client to forward the redeployment request
            var content = new StringContent(JsonSerializer.Serialize(request), Encoding.UTF8, "application/json");

            // Use the client to forward the redeployment request
            return await Client.PostAsync<StartRedeploymentResult>($"api/deployments/{deploymentId}/redeploy", content);
        }

        /// <summary>
        /// Gets the parameters associated with a deployment
        /// </summary>
        [HttpGet("deployments/{deploymentId}/parameters")]
        public async Task<IActionResult> GetParametersFileContent(string deploymentId)
        {
            var result = await Client.GetAsync<Dictionary<string, object>>(String.Format(Routes.GetDeploymentParameters, deploymentId));
            return result;
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
