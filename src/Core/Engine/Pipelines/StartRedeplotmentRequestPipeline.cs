using JenkinsNET.Models;
using JenkinsNET.Exceptions;
using MediatR;
using MediatR.Pipeline;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Modm.Deployments;
using Modm.Jenkins.Client;
using Modm.Engine.Notifications;
using static Microsoft.WindowsAzure.ResourceStack.Common.Utilities.FastActivator;

namespace Modm.Engine.Pipelines
{
    /// <summary>
    /// The pipeline
    /// </summary>
    public class StartRedeploymentRequestPipeline : Pipeline<StartRedeploymentRequest, StartRedeploymentResult>
    {
        public StartRedeploymentRequestPipeline(IMediator mediator) : base(mediator)
        {
        }
    }

    public static class StartRedeploymentRequestPipelineRegistration
    {
        public static MediatRServiceConfiguration AddStartRedeploymentRequestPipeline(this MediatRServiceConfiguration c)
        {
            // add sub pipeline which is a dependency for establishing the definition
            //c.AddCreateDeploymentDefinitionPipeline();

            c.RegisterServicesFromAssemblyContaining<StartRedeploymentRequestHandler>();

            //c.AddRequestPostProcessor<WriteDeploymentToDisk>();
            c.AddBehavior<SubmitRedeployment>();
            c.AddBehavior<DeleteParametersFile>();
            return c;
        }
    }

    public class StartRedeploymentRequestHandler : IRequestHandler<StartRedeploymentRequest, StartRedeploymentResult>
    {
        readonly IMediator mediator;

        public StartRedeploymentRequestHandler(IMediator mediator)
        {
            this.mediator = mediator;
        }

        public async Task<StartRedeploymentResult> Handle(StartRedeploymentRequest request, CancellationToken cancellationToken)
        {
            var returnedString = await GetString();

            return new StartRedeploymentResult
            {
                Deployment = new Deployment
                {
                    Definition = new DeploymentDefinition { DeploymentType = DeploymentType.Terraform, MainTemplatePath = returnedString }
                }
            };
        }

        private async Task<string> GetString()
        {
            return "test";
        }

        //private async Task<DeploymentDefinition> CreateDefinition(StartDeploymentRequest request, CancellationToken cancellationToken)
        //    => await mediator.Send<DeploymentDefinition>(new CreateDeploymentDefinition(request), cancellationToken);
    }

    public class SubmitRedeployment : IPipelineBehavior<StartRedeploymentRequest, StartRedeploymentResult>
    {
        private readonly JenkinsClientFactory clientFactory;
        private readonly IMediator mediator;
        private readonly ILogger<SubmitRedeployment> logger;

        public SubmitRedeployment(JenkinsClientFactory clientFactory, IMediator mediator, ILogger<SubmitRedeployment> logger)
        {
            this.clientFactory = clientFactory;
            this.mediator = mediator;
            this.logger = logger;
        }

        public async Task<StartRedeploymentResult> Handle(StartRedeploymentRequest request, RequestHandlerDelegate<StartRedeploymentResult> next, CancellationToken cancellationToken)
        {
            int x = 0;

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

        private async Task Publish(Deployment deployment, CancellationToken cancellationToken)
        {
            this.logger.LogInformation("Inside SubmitDeployment:Publish - publishing DeploymentStarted");
            await mediator.Publish(new DeploymentStarted
            {
                Id = deployment.Id,
                Name = deployment.Definition.DeploymentType
            }, cancellationToken);
        }

        private async Task<bool> TryToSubmit(Deployment deployment)
        {
            using var client = await clientFactory.Create();
            this.logger.LogInformation($"Prior to calling client.Build  - {DateTime.UtcNow}");


            var id = await client.Build(deployment.Definition.DeploymentType);

            // this was added in replace of commented section
            if (!id.HasValue)
            {
                this.logger.LogInformation("id does not have value in TryToSubmit");
                return false;
            }

            // If we get here, it means we have a valid build ID
            deployment.Id = id.Value;
            deployment.Status = DeploymentStatus.Undefined;
            this.logger.LogInformation($"The deployment.Id has a value of {deployment.Id}");
            return true;
        }

        private static void AddError(StartRedeploymentResult result, string error)
        {
            if (result.Errors == null)
            {
                result.Errors = new List<string>();
            }

            result.Errors.Add(error);
        }
    }

    public class DeleteParametersFile : IPipelineBehavior<StartRedeploymentRequest, StartRedeploymentResult>
    {
        private readonly ParametersFileFactory factory;
        private readonly IMediator mediator;
        private readonly ILogger<SubmitRedeployment> logger;

        public DeleteParametersFile(ParametersFileFactory factory, IMediator mediator, ILogger<SubmitRedeployment> logger)
        {
            this.factory = factory;
            this.mediator = mediator;
            this.logger = logger;
        }

        public async Task<StartRedeploymentResult> Handle(StartRedeploymentRequest request, RequestHandlerDelegate<StartRedeploymentResult> next, CancellationToken cancellationToken)
        {
         //   request.
         //   var file = factory.Create(definition.DeploymentType, definition.GetMainTemplateDirectoryName());
            var startDeploymentResult = new StartRedeploymentResult
            {
                Deployment = new Deployment { Id = 0 }
            };
      //      int y = 0;
            return startDeploymentResult;
        }
    }
}