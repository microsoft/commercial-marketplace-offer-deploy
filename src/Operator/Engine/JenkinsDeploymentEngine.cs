using System;
using Operator.Artifacts;

namespace Operator.Engine
{
    public class JenkinsDeploymentEngine : IDeploymentEngine
    {
        private readonly ArtifactsDownloader downloader;

        public JenkinsDeploymentEngine(ArtifactsDownloader downloader)
        {
            this.downloader = downloader;
        }


        /// <summary>
        /// starts a deployment
        /// </summary>
        /// <returns></returns>
        public async Task<int> StartAsync(string artifactsUri)
        {
            await downloader.DownloadAsync(new ArtifactsUri { Container = artifactsUri });

            return 1;
        }
    }
}

