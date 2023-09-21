using MediatR;
using Microsoft.Extensions.DependencyInjection;
using Modm.Artifacts;
using Modm.Deployments;

namespace Modm.Engine.Pipelines
{
    /// <summary>
    /// Child pipeline of <see cref="StartDeploymentRequestPipeline"/>
    /// </summary>
    public static class CreateDeploymentDefinitionPipeline
    {
        public static MediatRServiceConfiguration AddCreateDeploymentDefinitionPipeline(this MediatRServiceConfiguration c)
        {
            c.RegisterServicesFromAssemblyContaining<CreateDeploymentDefinitionHandler>();

            // start with behaviors order from bottom --> up
            // since we're going to handle the build up of the definition

            c.AddBehavior<CreateParametersFile>();
            c.AddBehavior<ReadManifestFile>();
            c.AddBehavior<DownloadAndExtractArtifactsFile>();

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
                Parameters = request.Parameters
            });
        }
    }


    /// <summary>
    /// first in the pipeline for creating a definition file
    /// </summary>
	public class DownloadAndExtractArtifactsFile : IPipelineBehavior<CreateDeploymentDefinition, DeploymentDefinition>
    {
        private readonly ArtifactsDownloader downloader;

        public DownloadAndExtractArtifactsFile(ArtifactsDownloader downloader)
        {
            this.downloader = downloader;
        }

        public async Task<DeploymentDefinition> Handle(CreateDeploymentDefinition request, RequestHandlerDelegate<DeploymentDefinition> next, CancellationToken cancellationToken)
        {
            var definition = await next();

            var artifactsFile = await downloader.DownloadAsync(definition.Source);
            artifactsFile.Extract();

            definition.WorkingDirectory = artifactsFile.ExtractedTo;

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

    #endregion
}