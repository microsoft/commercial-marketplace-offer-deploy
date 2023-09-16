using System;
using System.IO.Compression;
using System.Net.Http;
using Azure.Identity;
using Azure.Storage.Blobs;
using Azure.Storage.Blobs.Models;
using Microsoft.Extensions.Options;
using Modm.Deployments;
using Newtonsoft.Json;

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
        public Task<ArtifactsDescriptor> DownloadAsync(ArtifactsUri uri)
        {
            return DownloadAsync(uri, new ArtifactsDownloadOptions
            {
                SavePath = options.SavePath
            });
        }

        public async Task<ArtifactsDescriptor> DownloadAsync(ArtifactsUri uri, ArtifactsDownloadOptions options)
        {
            var httpResult = await client.GetAsync(uri);
            
            var resolvedSavePath = Environment.ExpandEnvironmentVariables(options.SavePath);
            var archiveFilePath = Path.GetFullPath(Path.Combine(resolvedSavePath, uri.FileName));

            using (var resultStream = await httpResult.Content.ReadAsStreamAsync())
            using (var fileStream = File.Create(archiveFilePath))
            {
                await resultStream.CopyToAsync(fileStream);
                await resultStream.FlushAsync();
            }

            ZipFile.ExtractToDirectory(archiveFilePath, resolvedSavePath, overwriteFiles: true);
            string manifestFilePath = Path.Combine(resolvedSavePath, Constants.ManifestFileName);


            if (File.Exists(manifestFilePath))
            {
                string manifestJson = File.ReadAllText(manifestFilePath);
                var definition = JsonConvert.DeserializeObject<DeploymentDefinition>(manifestJson);

                return new ArtifactsDescriptor
                {
                    FolderPath = resolvedSavePath,
                    Definition = definition
                };
            }
            else
            {
                throw new FileNotFoundException("manifest.json not found in the extracted files.");
            }
        }
    }
}

