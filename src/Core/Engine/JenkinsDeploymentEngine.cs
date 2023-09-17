using System;
using System.IO;
using System.Text.Json;
using JenkinsNET.Models;
using MediatR;
using Microsoft.Extensions.Configuration;
using Modm.Artifacts;
using Modm.Deployments;
using Modm.Engine.Jenkins;
using Modm.Engine.Jenkins.Client;
using Modm.Engine.Notifications;
using Modm.Extensions;

namespace Modm.Engine
{
    class JenkinsDeploymentEngine : IDeploymentEngine
    {
        readonly static JsonSerializerOptions DefaultSerializationOptions = new JsonSerializerOptions
        {
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
            WriteIndented = true
        };

        private readonly ArtifactsDownloader downloader;
        private readonly IJenkinsClient client;
        private readonly IConfiguration configuration;
        private readonly IMediator mediator;

        public JenkinsDeploymentEngine(ArtifactsDownloader downloader, IJenkinsClient client, IConfiguration configuration, IMediator mediator)
        {
            this.downloader = downloader;
            this.client = client;
            this.configuration = configuration;
            this.mediator = mediator;
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

        public async Task<Deployment> Get()
        {
            // TODO: do more than read from disk. We have to get the current status of the job
            return await Read();
        }


        /// <summary>
        /// starts a deployment
        /// </summary>
        /// <returns></returns>
        public async Task<StartDeploymentResult> Start(string artifactsUri)
        {
            var deployment = await Get();
            var descriptor = await downloader.DownloadAsync(new ArtifactsUri(artifactsUri));

            if (deployment.IsStartable)
            {
                deployment.Source = artifactsUri;
                deployment.Definition = descriptor.Definition;

                var buildResult = await client.Jobs.BuildAsync(descriptor.Definition.DeploymentType);
                var queueId = buildResult.GetQueueItemNumber().GetValueOrDefault(0);
                var queueItem = await client.Queue.GetItemAsync(queueId);

                deployment.Id = queueItem.Executable.Number.GetValueOrDefault(0);
                deployment.Status = DeploymentStatus.Running;

                await Write(deployment);

                await mediator.Publish(new DeploymentStarted
                {
                    Id = deployment.Id,
                    Name = deployment.Definition.DeploymentType
                });
            }

            return new StartDeploymentResult
            {
                Id = deployment.Id,
                Status = deployment.Status
            };
        }

        private async Task Write(Deployment deployment)
        {
            var json = JsonSerializer.Serialize(deployment, DefaultSerializationOptions);
            await File.WriteAllTextAsync(GetDeploymentFilePath(), json);
        }

        private async Task<Deployment> Read()
        {
            var path = GetDeploymentFilePath();

            if (!File.Exists(path))
            {
                return new Deployment
                {
                    Id = 0,
                    Status = DeploymentStatus.Undefined
                };
            }

            var json = await File.ReadAllTextAsync(GetDeploymentFilePath());
            var deployment = JsonSerializer.Deserialize<Deployment>(json, DefaultSerializationOptions);
            deployment.Status = await GetStatus(deployment);

            return deployment;
        }

        private string GetDeploymentFilePath()
        {
            return Path.GetFullPath(Path.Combine(configuration.GetHomeDirectory(), "deployment.json"));
        }

        private async Task<string> GetStatus(Deployment deployment)
        {
            var build = await client.Builds.GetAsync<JenkinsBuildBase>(deployment.Definition.DeploymentType, deployment.Id.ToString());
            return build.Result.ToLower();
        }
    }
}

