
using Microsoft.AspNetCore.Mvc;
using Modm;
using Modm.Engine;

namespace WebHost.Deployments
{
    [Route("api/[controller]")]
    [ApiController]
    public class ResourcesController : ControllerBase
    {
        [HttpGet]
        public Task<List<string>> GetAsync()
        {
            return Task.FromResult(new List<string>());
        }
    }
}
