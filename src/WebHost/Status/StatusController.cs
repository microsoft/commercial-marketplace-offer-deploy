using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using JenkinsNET;
using Modm.Engine.Jenkins;
using Modm.Engine;

// For more information on enabling MVC for empty projects, visit https://go.microsoft.com/fwlink/?LinkID=397860

namespace WebHost.Status
{
    [Route("api/[controller]")]
    [ApiController]
    public class StatusController : ControllerBase
    {
        private readonly IDeploymentEngine engine;

        public StatusController(IDeploymentEngine engine)
        {
            this.engine = engine;
        }

        [HttpGet]
        public async Task<EngineStatus> GetAsync()
        {
            return await this.engine.GetStatus();
        }
    }
}

