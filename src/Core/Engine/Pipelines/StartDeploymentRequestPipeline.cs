﻿using JenkinsNET.Models;
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
    }

    // #3
    public class WriteDeploymentToDisk : IRequestPostProcessor<StartDeploymentRequest, StartDeploymentResult>
    {
        private readonly DeploymentFile deploymentFile;
        private readonly AuditFile auditFile;
        private readonly ILogger<WriteDeploymentToDisk> logger;

        public WriteDeploymentToDisk(DeploymentFile deploymentFile, AuditFile auditFile, ILogger<WriteDeploymentToDisk> logger)
        {
            this.deploymentFile = deploymentFile;
            this.auditFile = auditFile;
            this.logger = logger;
        }
        public async Task Process(
            StartDeploymentRequest request,
            StartDeploymentResult response,
            CancellationToken cancellationToken)
        {
            this.logger.LogInformation("Inside WriteDeploymentToDisk:Process");

            var deployment = await this.deploymentFile.ReadAsync(cancellationToken);
            await deploymentFile.WriteAsync(deployment, cancellationToken);

            var auditRecords = await this.auditFile.ReadAsync(cancellationToken);
            var auditRecord = new AuditRecord();
            auditRecord.AdditionalData.Add("WriteDeploymentToDisk:Process", response.Deployment);
            auditRecords.Add(auditRecord);
            await this.auditFile.WriteAsync(auditRecords, cancellationToken);
        }
    }

    #endregion
}

