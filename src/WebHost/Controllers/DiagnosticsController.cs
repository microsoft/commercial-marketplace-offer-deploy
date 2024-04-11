using Microsoft.AspNetCore.Mvc;
using Modm.Diagnostics;
using Modm.Engine;
using Modm.Extensions;

namespace WebHost.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class DiagnosticsController : ControllerBase
    {
        private readonly string stripMarker = "-----------------";
        private readonly IDeploymentEngine engine;

        public DiagnosticsController(IDeploymentEngine engine)
        {
            this.engine = engine;
        }

        public async Task<IResult> Get()
        {
            string logsContent = await engine.GetLogs();

            return Results.Json(new GetDiagnosticsResponse
            {
                DeploymentEngine = logsContent
                    .SubstringAfterMarker(this.stripMarker)
                    .Replace("jenkins", "***")
            }); 
        }
    }
}