using System;
using Microsoft.AspNetCore.DataProtection.KeyManagement;
using Microsoft.AspNetCore.Mvc;

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
    }
}