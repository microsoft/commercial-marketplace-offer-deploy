using Microsoft.AspNetCore.Mvc;
using Modm.Engine;

// For more information on enabling MVC for empty projects, visit https://go.microsoft.com/fwlink/?LinkID=397860

namespace WebHost.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class StatusController : ControllerBase
    {
        private readonly IDeploymentEngine engine;
        private readonly JenkinsReadinessService readinessService;

        public StatusController(IDeploymentEngine engine, JenkinsReadinessService readinessService)
        {
            this.engine = engine;
            this.readinessService = readinessService;
        }

        [HttpGet]
        public EngineInfo Get()
        {
            //return await this.engine.GetInfo();
            return this.readinessService.GetEngineInfo();
        }
    }
}

