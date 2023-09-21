using System;
using FluentValidation;
using Microsoft.AspNetCore.Mvc;
using Modm.Deployments;
using Modm.Diagnostics;
using Modm.Engine;

namespace WebHost.Controllers
{

    [Route("api/[controller]")]
    [ApiController]
    public class DiagnosticsController : ControllerBase
    {
        private readonly IDeploymentEngine engine;

        public DiagnosticsController(IDeploymentEngine engine)
        {
            this.engine = engine;
        }

        public async Task<IResult> Get()
        {
            return Results.Json(new GetDiagnosticsResponse
            {
                DeploymentEngine = await engine.GetLogs()
            });
        }
    }
}

