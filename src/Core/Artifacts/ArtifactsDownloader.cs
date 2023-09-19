using System;
using System.IO.Compression;
using System.Net.Http;
using Azure.Identity;
using Azure.Storage.Blobs;
using Azure.Storage.Blobs.Models;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Options;
using Modm.Deployments;
using Modm.Extensions;
using Newtonsoft.Json;

namespace Modm.Artifacts
{
	public class ArtifactsDownloader
	{
        private readonly HttpClient client;
        private readonly IConfiguration configuration;

        public ArtifactsDownloader(HttpClient client, IConfiguration configuration)
		{
            this.client = client;
            this.configuration = configuration;
        }

        /// <summary>
        /// save the artifacts from uri to the configured save path in appsettings
        /// </summary>
        /// <param name="uri"></param>
        /// <returns></returns>
        public Task<ArtifactsDescriptor> DownloadAsync(ArtifactsUri uri)
        {
            return DownloadAsync(uri, new ArtifactsDownloadOptions
            {
                SavePath = configuration.GetHomeDirectory()
            });
        }

        public async Task<ArtifactsDescriptor> DownloadAsync(ArtifactsUri uri, ArtifactsDownloadOptions options)
        {
            var httpResult = await client.GetAsync(uri);
            
            var archiveFilePath = Path.GetFullPath(Path.Combine(options.SavePath, uri.FileName));

            using (var resultStream = await httpResult.Content.ReadAsStreamAsync())
            using (var fileStream = File.Create(archiveFilePath))
            {
                await resultStream.CopyToAsync(fileStream);
                await resultStream.FlushAsync();
            }

            var artifactsFile = new ArtifactsFile(archiveFilePath);
            var definition = await artifactsFile.ReadManifestFile();

            return new ArtifactsDescriptor
            {
                FolderPath = artifactsFile.ExtractedTo,
                Definition = definition
            };
        }
    }
}

