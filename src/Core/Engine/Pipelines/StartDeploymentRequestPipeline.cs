using JenkinsNET.Models;
using JenkinsNET.Exceptions;
using MediatR;
using MediatR.Pipeline;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Modm.Deployments;
using Modm.Jenkins.Client;
using Modm.Engine.Notifications;

namespace Modm.Engine.Pipelines
{
    /// <summary>
    /// The pipeline
    /// </summary>
    public class StartDeploymentRequestPipeline : Pipeline<StartDeploymentRequest, StartDeploymentResult>
    {
        public StartDeploymentRequestPipeline(IMediator mediator) : base(mediator)
        {
        }
    }

    public static class StartDeploymentRequestPipelineRegistration
    {
        public static MediatRServiceConfiguration AddStartDeploymentRequestPipeline(this MediatRServiceConfiguration c)
        {
            // add sub pipeline which is a dependency for establishing the definition
            c.AddCreateDeploymentDefinitionPipeline();

            c.RegisterServicesFromAssemblyContaining<StartDeploymentRequestHandler>();

            c.AddRequestPostProcessor<WriteDeploymentToDisk>();
            c.AddBehavior<SubmitDeployment>();
            c.AddBehavior<ReadDeploymentFromRepository>();
            return c;
        }
    }

    #region Pipeline

    /// <summary>
    /// Starts the deployment request
    /// </summary>
    public class StartDeploymentRequestHandler : IRequestHandler<StartDeploymentRequest, StartDeploymentResult>
    {
        readonly IMediator mediator;

        public StartDeploymentRequestHandler(IMediator mediator)
        {
            this.mediator = mediator;
        }

        public async Task<StartDeploymentResult> Handle(StartDeploymentRequest request, CancellationToken cancellationToken)
        {
            var definition = await CreateDefinition(request, cancellationToken);

            return new StartDeploymentResult
            {
                Deployment = new Deployment
                {
                    Definition = definition
                }
            };
        }

        private async Task<DeploymentDefinition> CreateDefinition(StartDeploymentRequest request, CancellationToken cancellationToken)
            => await mediator.Send<DeploymentDefinition>(new CreateDeploymentDefinition(request), cancellationToken);
    }

    // #1
    public class ReadDeploymentFromRepository : IPipelineBehavior<StartDeploymentRequest, StartDeploymentResult>
    {
        private readonly IDeploymentRepository repository;
        public ReadDeploymentFromRepository(IDeploymentRepository repository) => this.repository = repository;

        public async Task<StartDeploymentResult> Handle(StartDeploymentRequest request, RequestHandlerDelegate<StartDeploymentResult> next, CancellationToken cancellationToken)
        {
            var result = await next();
            result.Deployment = await repository.Get(cancellationToken);

            return result;
        }
    }

    // #2
    public class SubmitDeployment : IPipelineBehavior<StartDeploymentRequest, StartDeploymentResult>
    {
        private readonly JenkinsClientFactory clientFactory;
        private readonly IMediator mediator;
        private readonly ILogger<SubmitDeployment> logger;

        public SubmitDeployment(JenkinsClientFactory clientFactory, IMediator mediator, ILogger<SubmitDeployment> logger)
        {
            this.clientFactory = clientFactory;
            this.mediator = mediator;
            this.logger = logger;
        }

        public async Task<StartDeploymentResult> Handle(StartDeploymentRequest request, RequestHandlerDelegate<StartDeploymentResult> next, CancellationToken cancellationToken)
        {
            var result = await next();
            result.Errors ??= new List<string>();

            var deployment = result.Deployment;
            
            if (!deployment.IsStartable)
            {
                deployment.Id = -1;
                AddError(result, "Deployment is not startable");
                return result;
            }

            try
            {
                if (await TryToSubmit(deployment))
                {
                    await Publish(deployment, cancellationToken);
                }
            }
            catch (Exception ex)
            {
                AddError(result, ex.Message);
                logger.LogError(ex, "Failure to submit to jenkins");
            }

            result.Deployment = deployment;

            return result;
        }

        private static void AddError(StartDeploymentResult result, string error)
        {
            if (result.Errors == null)
            {
                result.Errors = new List<string>();
            }

            result.Errors.Add(error);
        }


        private async Task Publish(Deployment deployment, CancellationToken cancellationToken)
        {
            await mediator.Publish(new DeploymentStarted
            {
                Id = deployment.Id,
                Name = deployment.Definition.DeploymentType
            }, cancellationToken);
        }

        private async Task<bool> TryToSubmit(Deployment deployment)
        {
            using var client = await clientFactory.Create();
            var (id, status) = await client.Build(deployment.Definition.DeploymentType);
            deployment.Status = status;

            if (id.HasValue)
            {
                deployment.Id = id.Value;
                return true;
            }
            return false;
        }
    }

    // #3
    public class WriteDeploymentToDisk : IRequestPostProcessor<StartDeploymentRequest, StartDeploymentResult>
    {
        private readonly DeploymentFile file;

        public WriteDeploymentToDisk(DeploymentFile file) => this.file = file;

        public async Task Process(
            StartDeploymentRequest request,
            StartDeploymentResult response,
            CancellationToken cancellationToken)
        {
            await file.Write(response.Deployment, cancellationToken);
        }
    }

    #endregion
}

