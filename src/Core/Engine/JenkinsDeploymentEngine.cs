using System;
using JenkinsNET;
using JenkinsNET.Models;
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

            // token = await this.apiTokenProvider.GetAsync();

            //var jobUrlWithToken = $"http://localhost:8083/job/modmserviceprincipal/job/modmserviceprincipal/build?token={token}";

            //var jenkinsClient = new JenkinsClient(jenkinsUrl);
            //bool isJobTriggered = await jenkinsClient.Jobs.TriggerRemoteJobAsync(jobUrlWithToken);

            //var httpClient = new HttpClient();
            // Send an HTTP POST request to trigger the job
            //HttpResponseMessage response = await httpClient.PostAsync(jobUrlWithToken, null);


            // TODO: use result.GetQueueItemNumber() and whatever else to have the backend start to poll for the information
            //string jobName = "modmserviceprincipal/modmserviceprincipal";
            var jobInfo =  client.Jobs.Get<JenkinsJobBase>("modmserviceprincipal");
            var result = await client.Jobs.BuildAsync("modmserviceprincipal");
            //client.Jobs.Get


            //var result = await client.Jobs.BuildAsync(descriptor.Definition.DeploymentType);
            // client.
            // this.apiTokenProvider.

            // TODO: fill out results. this is just stubbed out only
            return new StartDeploymentResult();
        }
    }
}

