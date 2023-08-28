﻿using System;
using JenkinsNET;
using Modm.Artifacts;
using Modm.Engine.Jenkins;

namespace Modm.Engine
{
    public class JenkinsDeploymentEngine : IDeploymentEngine
    {
        private readonly ArtifactsDownloader downloader;
        private readonly ApiTokenProvider apiTokenProvider;
        private readonly IJenkinsClient client;

        public JenkinsDeploymentEngine(ArtifactsDownloader downloader, ApiTokenProvider apiTokenProvider, IJenkinsClient client)
        {
            this.downloader = downloader;
            this.apiTokenProvider = apiTokenProvider;
            this.client = client;
        }


        /// <summary>
        /// starts a deployment
        /// </summary>
        /// <returns></returns>
        public async Task<StartDeploymentResult> StartAsync(string artifactsUri)
        {
            var descriptor = await downloader.DownloadAsync(new ArtifactsUri(artifactsUri));

            // TODO: use result.GetQueueItemNumber() and whatever else to have the backend start to poll for the information
            var result = await client.Jobs.BuildAsync(descriptor.Definition.DeploymentType);
           
            // TODO: fill out results. this is just stubbed out only
            return new StartDeploymentResult();
        }
    }
}
