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
