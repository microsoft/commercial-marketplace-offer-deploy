using System;
using Microsoft.AspNetCore.DataProtection.KeyManagement;
using Microsoft.AspNetCore.Mvc;
using Modm.Engine;
using ClientApp.Backend;
using Azure.ResourceManager;
using Modm.Azure;
using MediatR;
using ClientApp.Commands;

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

        [HttpGet("deployments")]
        public async Task<IActionResult> GetDeployments()
        {
            // TODO: implement proxy to backend with Http Cient
            return Task.FromResult(Results.Ok(null));
        }

        [HttpGet("status")]
        public async Task<IActionResult> GetStatus()
        {
            // TODO: implement proxy to backend with Http Cient
            return Task.FromResult(Results.Ok(EngineInfo.Default()));
        }
    }
}
