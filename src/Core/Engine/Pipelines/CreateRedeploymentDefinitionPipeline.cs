using FluentValidation;
using MediatR;
using MediatR.Pipeline;
using Microsoft.Extensions.DependencyInjection;
using Modm.Packaging;
using Modm.Deployments;
using Microsoft.Extensions.Logging;

namespace Modm.Engine.Pipelines
{
    /// <summary>
    /// Child pipeline of <see cref="StartDeploymentRequestPipeline"/>
    /// </summary>
    public static class CreateRedeploymentDefinitionPipeline
    {
        public static MediatRServiceConfiguration AddCreateRedeploymentDefinitionPipeline(this MediatRServiceConfiguration c)
        {
            // start with behaviors order from bottom --> up
            // since we're going to handle the build up of the definition
   
            c.AddBehavior<CreateParametersFile>();
            c.AddBehavior<ReadManifestFile>();
            c.AddRequestPostProcessor<WriteToDisk>();

            c.RegisterServicesFromAssemblyContaining<CreateDeploymentDefinitionHandler>();

   
            return c;
        }
    }

    #region Pipeline

    /// <summary>
    /// Starts the pipeline of the definition creation
    /// </summary>
    public class CreateRedeploymentDefinitionHandler : IRequestHandler<CreateRedeploymentDefinition, DeploymentDefinition>
    {
        public Task<DeploymentDefinition> Handle(CreateRedeploymentDefinition request, CancellationToken cancellationToken)
        {
            return Task.FromResult(new DeploymentDefinition
            {
                Source = request.GetUri(),
                InstallerPackageHash = request.PackageHash,
                Parameters = request.Parameters
            });
        }
    }

    public class ReadDeploymentFile : IPipelineBehavior<CreateRedeploymentDefinition, DeploymentDefinition>
    {
        private readonly DeploymentFile deploymentFile;

        public ReadDeploymentFile(DeploymentFile deploymentFile)
        {
            this.deploymentFile = deploymentFile;

        }

        public async Task<DeploymentDefinition> Handle(CreateRedeploymentDefinition request, RequestHandlerDelegate<DeploymentDefinition> next, CancellationToken cancellationToken)
        {
            var result = await next();

            var deployment = await this.deploymentFile.ReadAsync(cancellationToken);
           // result.

            return result;
        }
    }

    

    // #2
    public class RereadManifestFile : IPipelineBehavior<CreateRedeploymentDefinition, DeploymentDefinition>
    {
        public async Task<DeploymentDefinition> Handle(CreateRedeploymentDefinition request, RequestHandlerDelegate<DeploymentDefinition> next, CancellationToken cancellationToken)
        {
            var definition = await next();

            var manifest = await ManifestFile.Read(definition.WorkingDirectory);

            definition.MainTemplatePath = manifest.MainTemplate;
            definition.DeploymentType = manifest.DeploymentType;

            return definition;
        }
    }

    // #3
    public class RecreateParametersFile : IPipelineBehavior<CreateRedeploymentDefinition, DeploymentDefinition>
    {
        private readonly ParametersFileFactory factory;

        public RecreateParametersFile(ParametersFileFactory parametersFileFactory)
        {
            this.factory = parametersFileFactory;
        }

        public async Task<DeploymentDefinition> Handle(CreateRedeploymentDefinition request, RequestHandlerDelegate<DeploymentDefinition> next, CancellationToken cancellationToken)
        {
            var definition = await next();
            var file = factory.Create(definition.DeploymentType, definition.GetMainTemplateDirectoryName());

            // the file must always have at least an empty object
            await file.Write(request.Parameters ?? new Dictionary<string, object>());
            definition.ParametersFilePath = file.FullPath;

            return definition;
        }
    }

    // #4
    public class RewriteToDisk : IRequestPostProcessor<CreateRedeploymentDefinition, DeploymentDefinition>
    {
        private readonly DeploymentFile deploymentFile;
        private readonly AuditFile auditFile;
        private ILogger<WriteToDisk> logger;

        public RewriteToDisk(DeploymentFile deploymentFile, AuditFile auditFile, ILogger<WriteToDisk> logger)
        {
            this.deploymentFile = deploymentFile;
            this.auditFile = auditFile;
            this.logger = logger;
        }

        public async Task Process(CreateRedeploymentDefinition request, DeploymentDefinition response, CancellationToken cancellationToken)
        {
            this.logger.LogInformation("Inside RewriteToDisk of CreateRedeploymentPipeline");

            var deployment = new Deployment
            {
                Definition = response,
                Id = request.DeploymentId,
                Timestamp = DateTimeOffset.UtcNow,
                Status = DeploymentStatus.Undefined
            };

            await deploymentFile.WriteAsync(deployment, cancellationToken);
            this.logger.LogInformation("Wrote Deployment to deployment file");

            var auditRecord = new AuditRecord();
            auditRecord.AdditionalData.Add("createDeploymentPipeline", deployment);

            await this.auditFile.WriteAsync(new List<AuditRecord>() { auditRecord }, cancellationToken);
        } 
    }

    #endregion
}