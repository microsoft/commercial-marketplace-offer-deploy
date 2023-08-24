using System;
using System.IO.Compression;
using System.Net.Http;
using Azure.Identity;
using Azure.Storage.Blobs;
using Azure.Storage.Blobs.Models;
using Microsoft.Extensions.Options;

namespace Modm.Artifacts
{
	public class ArtifactsDownloader
	{
        private readonly HttpClient client;
        private readonly ArtifactsDownloadOptions options;

        public ArtifactsDownloader(HttpClient client, IOptions<ArtifactsDownloadOptions> options)
		{
            this.client = client;
            this.options = options.Value;
        }

        /// <summary>
        /// save the artifacts from uri to the configured save path in appsettings
        /// </summary>
        /// <param name="uri"></param>
        /// <returns></returns>
        public Task DownloadAsync(ArtifactsUri uri)
        {
            return DownloadAsync(uri, new ArtifactsDownloadOptions
            {
                SavePath = options.SavePath
            });
        }

        public async Task DownloadAsync(ArtifactsUri uri, ArtifactsDownloadOptions options)
        {
            var httpResult = await client.GetAsync(uri);
            var archiveFilePath = Path.Combine(options.SavePath, uri.FileName);

            using (var resultStream = await httpResult.Content.ReadAsStreamAsync())
            using (var fileStream = File.Create(archiveFilePath))
            {
                await resultStream.CopyToAsync(fileStream);
                await resultStream.FlushAsync();
            }

            ZipFile.ExtractToDirectory(archiveFilePath, options.SavePath, overwriteFiles: true);
        }
    }
}

