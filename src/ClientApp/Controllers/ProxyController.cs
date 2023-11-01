using System;
using Microsoft.AspNetCore.DataProtection.KeyManagement;
using Microsoft.AspNetCore.Mvc;
using Modm.Engine;

namespace Modm.ClientApp.Controllers
{
    [Route("api")]
    [ApiController]
    public class ProxyController : ControllerBase
    {
        public const string BackendUrlSettingName = "BackendUrl";
        private readonly HttpClient client;
        private readonly IConfiguration configuration;

        public ProxyController(HttpClient client, IConfiguration configuration)
        {
            this.client = client;
            this.configuration = configuration;
        }

        [HttpGet("deployments")]
        public Task<IResult> GetDeployments()
        {
            // TODO: implement proxy to backend with Http Cient
            return Task.FromResult(Results.Ok(null));
        }

        [HttpGet("status")]
        public Task<IResult> GetStatus()
        {
            // TODO: implement proxy to backend with Http Cient
            return Task.FromResult(Results.Ok(EngineInfo.Default()));
        }
    }
}