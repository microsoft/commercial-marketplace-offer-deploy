using Microsoft.Extensions.Configuration;
using Modm.Extensions;

namespace Modm.Artifacts
{
    public class ArtifactsDownloader : IArtifactsDownloader
    {
        private readonly HttpClient client;
        private readonly IConfiguration configuration;
        private readonly ArtifactsFileFactory artifactsFileFactory;

        public ArtifactsDownloader(HttpClient client, IConfiguration configuration, ArtifactsFileFactory artifactsFileFactory)
        {
            this.client = client;
            this.configuration = configuration;
            this.artifactsFileFactory = artifactsFileFactory;
        }

        /// <summary>
        /// save the artifacts from uri to the configured save path in appsettings
        /// </summary>
        /// <param name="uri"></param>
        /// <returns></returns>
        public Task<ArtifactsFile> DownloadAsync(ArtifactsUri uri)
        {
            return DownloadAsync(uri, new ArtifactsDownloadOptions
            {
                SavePath = configuration.GetHomeDirectory()
            });
        }

        public async Task<ArtifactsFile> DownloadAsync(ArtifactsUri uri, ArtifactsDownloadOptions options)
        {
            var httpResult = await client.GetAsync(uri);
            var artifactsFile = await DownloadFile(httpResult, options);

            return artifactsFile;
        }

        private async Task<ArtifactsFile> DownloadFile(HttpResponseMessage httpResult, ArtifactsDownloadOptions options)
        {
            var archiveFilePath = Path.GetFullPath(Path.Combine(options.SavePath, ArtifactsFile.FileName));

            using var resultStream = await httpResult.Content.ReadAsStreamAsync();
            using var fileStream = File.Create(archiveFilePath);

            await resultStream.CopyToAsync(fileStream);
            await resultStream.FlushAsync();

            return artifactsFileFactory.Create(archiveFilePath);
        }
    }
}

