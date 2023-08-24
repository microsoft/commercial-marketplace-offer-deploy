using System;
using Modm.Artifacts;

namespace Modm.Engine
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
            await downloader.DownloadAsync(new ArtifactsUri(artifactsUri));

            return 1;
        }
    }
}

