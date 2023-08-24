using FluentValidation;
using Microsoft.AspNetCore.Mvc;
using Modm;
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

        /// <summary>
        /// Starts a deployment by submitting to the Operator's deployment engine
        /// </summary>
        [HttpPost]
        public async Task<IResult> PostAsync([FromBody] CreateDeploymentRequest request)
        {
            var validationResult = await validator.ValidateAsync(request);

            if (!validationResult.IsValid)
            {
                return Results.ValidationProblem(validationResult.ToDictionary());
            }

            var id = await engine.StartAsync(request.ArtifactsUri);


            // TODO: get resulting object from repo. (underlying "repository" for deployments is going to be the deployment engine)
            // the underlying engine will be jenkins. the wrapper for jenkins will use the jenkins job id as the id of the deployment attempt
            return Results.Created("/deployments", new CreateDeploymentResponse
            {
                Id = id
            });
        }
    }
}
