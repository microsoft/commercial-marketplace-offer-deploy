using System;
using System.IO;
using System.Text.Json;
using JenkinsNET.Models;
using Microsoft.Extensions.Configuration;
using Modm.Artifacts;
using Modm.Deployments;
using Modm.Engine.Jenkins;
using Modm.Engine.Jenkins.Client;
using Modm.Extensions;

namespace Modm.Engine
{
    class JenkinsDeploymentEngine : IDeploymentEngine
    {
        private readonly ArtifactsDownloader downloader;
        private readonly IJenkinsClient client;
        private readonly IConfiguration configuration;

        public JenkinsDeploymentEngine(ArtifactsDownloader downloader, IJenkinsClient client, IConfiguration configuration)
        {
            this.downloader = downloader;
            this.client = client;
            this.configuration = configuration;
        }

        public async Task<EngineInfo> GetInfo()
        {
            var info = await client.GetInfo();
            var node = await client.GetBuiltInNode();

            return new EngineInfo
            {
                EngineType = EngineType.Jenkins,
                Version = info.Version,
                IsHealthy = !node.Offline
            };
        }

        public Task<EngineStatus> GetStatus()
        {
            return Task.FromResult(new EngineStatus { });
        }


        /// <summary>
        /// starts a deployment
        /// </summary>
        /// <returns></returns>
        public async Task<StartDeploymentResult> Start(string artifactsUri)
        {
            var deployment = await Get();
            var descriptor = await downloader.DownloadAsync(new ArtifactsUri(artifactsUri));

            if (deployment.Status == DeploymentStatus.Undefined)
            {
                deployment.Source = artifactsUri;
                deployment.Definition = descriptor.Definition;

                var result = await client.Jobs.BuildAsync(descriptor.Definition.DeploymentType);

                deployment.Id = result.GetQueueItemNumber().GetValueOrDefault(0);
                deployment.Status = DeploymentStatus.Running;

                await Save(deployment);
            }

            return new StartDeploymentResult
            {
                Id = deployment.Id
            };
        }

        private async Task Save(Deployment deployment)
        {
            var json = JsonSerializer.Serialize(deployment);
            await File.WriteAllTextAsync(GetDeploymentFilePath(), json);
        }

        private async Task<Deployment> Get()
        {
            var path = GetDeploymentFilePath();

            if (!File.Exists(path))
            {
                return new Deployment
                {
                    Status = DeploymentStatus.Undefined
                };
            }

            var json = await File.ReadAllTextAsync(GetDeploymentFilePath());
            return JsonSerializer.Deserialize<Deployment>(json);
        }

        private string GetDeploymentFilePath()
        {
            return Path.GetFullPath(Path.Combine(configuration.GetHomeDirectory(), "deployment.json"));
        }
    }
}

