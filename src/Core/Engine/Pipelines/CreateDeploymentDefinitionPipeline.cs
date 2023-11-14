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
    public static class CreateDeploymentDefinitionPipeline
    {
        public static MediatRServiceConfiguration AddCreateDeploymentDefinitionPipeline(this MediatRServiceConfiguration c)
        {
            // start with behaviors order from bottom --> up
            // since we're going to handle the build up of the definition
   
            c.AddBehavior<CreateParametersFile>();
            c.AddBehavior<ReadManifestFile>();
            c.AddBehavior<DownloadAndExtractInstallerPackage>();
            c.AddRequestPostProcessor<WriteToDisk>();

            c.RegisterServicesFromAssemblyContaining<CreateDeploymentDefinitionHandler>();

   
            return c;
        }
    }

    #region Pipeline

    /// <summary>
    /// Starts the pipeline of the definition creation
    /// </summary>
    public class CreateDeploymentDefinitionHandler : IRequestHandler<CreateDeploymentDefinition, DeploymentDefinition>
    {
        public Task<DeploymentDefinition> Handle(CreateDeploymentDefinition request, CancellationToken cancellationToken)
        {
            return Task.FromResult(new DeploymentDefinition
            {
                Source = request.GetUri(),
                InstallerPackageHash = request.PackageHash,
                Parameters = request.Parameters
            });
        }
    }

    // #1
    /// <summary>
    /// first in the pipeline for creating a definition file
    /// </summary>
	public class DownloadAndExtractInstallerPackage : IPipelineBehavior<CreateDeploymentDefinition, DeploymentDefinition>
    {
        private readonly IPackageDownloader downloader;
        private readonly IServiceProvider serviceProvider;
        //private readonly IValidator<PackageFile> validator;

        //public DownloadAndExtractInstallerPackage(IPackageDownloader downloader, IValidator<PackageFile> validator)
        //{
        //    this.downloader = downloader;
        //    this.validator = validator;
        //}

        public DownloadAndExtractInstallerPackage(IPackageDownloader downloader, IServiceProvider serviceProvider)
        {
            this.downloader = downloader;
            this.serviceProvider = serviceProvider;
        }

        public async Task<DeploymentDefinition> Handle(CreateDeploymentDefinition request, RequestHandlerDelegate<DeploymentDefinition> next, CancellationToken cancellationToken)
        {
            var definition = await next();
            definition.InstallerPackageHash = request.PackageHash;

            var file = await downloader.DownloadAsync(definition.Source);

            var context = new ValidationContext<PackageFile>(file);
            context.RootContextData[PackageFile.HashAttributeName] = request.PackageHash;

            using (var scope = this.serviceProvider.CreateScope())
            {
                var validator = scope.ServiceProvider.GetRequiredService<IValidator<PackageFile>>();
                var validationResult = validator.Validate(context);

                if (!validationResult.IsValid)
                {
                    throw new ValidationException("Error handling installer package extraction", validationResult.Errors);
                }
            }
            
            file.Extract();
            definition.WorkingDirectory = file.ExtractedTo;

            return definition;
        }
    }

    // #2
    public class ReadManifestFile : IPipelineBehavior<CreateDeploymentDefinition, DeploymentDefinition>
    {
        public async Task<DeploymentDefinition> Handle(CreateDeploymentDefinition request, RequestHandlerDelegate<DeploymentDefinition> next, CancellationToken cancellationToken)
        {
            var definition = await next();

            var manifest = await ManifestFile.Read(definition.WorkingDirectory);

            definition.MainTemplatePath = manifest.MainTemplate;
            definition.DeploymentType = manifest.DeploymentType;

            return definition;
        }
    }

    // #3
    public class CreateParametersFile : IPipelineBehavior<CreateDeploymentDefinition, DeploymentDefinition>
    {
        public async Task<DeploymentDefinition> Handle(CreateDeploymentDefinition request, RequestHandlerDelegate<DeploymentDefinition> next, CancellationToken cancellationToken)
        {
            var definition = await next();

            var factory = new ParametersFileFactory();
            var file = factory.Create(definition.DeploymentType, definition.GetMainTemplateDirectoryName());

            // the file must always have at least an empty object
            await file.Write(request.Parameters ?? new Dictionary<string, object>());
            definition.ParametersFilePath = file.FullPath;

            return definition;
        }
    }

    // #4
    public class WriteToDisk : IRequestPostProcessor<CreateDeploymentDefinition, DeploymentDefinition>
    {
        private readonly DeploymentFile file;
        private ILogger<WriteToDisk> logger;

        public WriteToDisk(DeploymentFile file, ILogger<WriteToDisk> logger)
        {
            this.file = file;
            this.logger = logger;
        }

        public async Task Process(CreateDeploymentDefinition request, DeploymentDefinition response, CancellationToken cancellationToken)
        {
            this.logger.LogInformation("Inside WriteToDisk of CreateDeploymentPipeline");

            await file.Write(new Deployment
            {
                Definition = response,
                Id = 0,
                Timestamp = DateTime.UtcNow,
                Status = DeploymentStatus.Undefined
            }, cancellationToken);

            this.logger.LogInformation("Wrote Deployment to file");
        } 
    }

    #endregion
}