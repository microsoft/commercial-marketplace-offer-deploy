using System.IO;
using FluentValidation;
using Microsoft.AspNetCore.Mvc;
using Modm.Deployments;
using Modm.Engine;

namespace WebHost.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class DeploymentsController : ControllerBase
    {
        private readonly IValidator<StartDeploymentRequest> validator;
        private readonly IDeploymentEngine engine;

        public DeploymentsController(IValidator<StartDeploymentRequest> validator, IDeploymentEngine engine)
        {
            this.validator = validator;
            this.engine = engine;
        }

        public async Task<IResult> Get()
        {
            return Results.Json(new GetDeploymentResponse
            {
                Deployment = await engine.Get()
            });
        }

        [HttpGet("{deploymentId}/parameters")]
        public async Task<IActionResult> GetParametersFileContent(string deploymentId)
        {
            try
            {
                Deployment deployment = await engine.Get();

                // Determine the path to the parameters file
                string parametersFilePath = deployment?.Definition?.ParametersFilePath;

                if (!System.IO.File.Exists(parametersFilePath))
                {
                    throw new FileNotFoundException();
                }

                string content = await System.IO.File.ReadAllTextAsync(parametersFilePath);

                return new ContentResult
                {
                    ContentType = "application/json",
                    Content = content,
                    StatusCode = 200
                };
            }
            catch (FileNotFoundException)
            {
                return NotFound("Parameters file not found.");
            }
            catch (Exception ex)
            {
                // Log the exception details
                // Return a generic error message to the client
                return StatusCode(500, "An error occurred while processing your request.");
            }
        }

        /// <summary>
        /// Creates a deployment by submitting to the deployment engine
        /// </summary>
        [HttpPost]
        public async Task<IResult> PostAsync([FromBody] StartDeploymentRequest request, CancellationToken cancellationToken)
        {
            var validationResult = await validator.ValidateAsync(request, cancellationToken);

            if (!validationResult.IsValid)
            {
                return Results.ValidationProblem(validationResult.ToDictionary());
            }

            var result = await engine.Start(request, cancellationToken);
            return Results.Created("/deployments", result);
        }
    }
}
