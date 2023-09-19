using FluentValidation;
using Microsoft.AspNetCore.Mvc;
using Modm.Deployments;
using Modm.Engine;

namespace WebHost.Deployments
{
    [Route("api/[controller]")]
    [ApiController]
    public class DeploymentsController : ControllerBase
    {
        private readonly IValidator<CreateDeploymentRequest> validator;
        private readonly IDeploymentEngine engine;

        public DeploymentsController(IValidator<CreateDeploymentRequest> validator, IDeploymentEngine engine)
        {
            this.validator = validator;
            this.engine = engine;
        }

        public async Task<IResult> Get()
        {
            return Results.Created("/deployments", new GetDeploymentResponse
            {
                Deployment = await engine.Get()
            });
        }

        /// <summary>
        /// Creates a deployment by submitting to the deployment engine
        /// </summary>
        [HttpPost]
        public async Task<IResult> PostAsync([FromBody] CreateDeploymentRequest request)
        {
            var validationResult = await validator.ValidateAsync(request);

            if (!validationResult.IsValid)
            {
                return Results.ValidationProblem(validationResult.ToDictionary());
            }

            var result = await engine.Start(request.ArtifactsUri, request.Parameters);

            return Results.Created("/deployments", new CreateDeploymentResponse
            {
                Id = result.Id,
                Status = result.Status
            });
        }
    }
}
