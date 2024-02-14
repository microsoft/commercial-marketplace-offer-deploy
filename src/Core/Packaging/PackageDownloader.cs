// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using Microsoft.Extensions.Configuration;
using Modm.Extensions;

namespace Modm.Packaging
{
    /// <summary>
    /// The installer package downloader
    /// </summary>
    public class PackageDownloader : IPackageDownloader
    {
        private readonly HttpClient client;
        private readonly IConfiguration configuration;
        private readonly PackageFileFactory factory;

        public PackageDownloader(HttpClient client, IConfiguration configuration, PackageFileFactory factory)
        {
            this.client = client;
            this.configuration = configuration;
            this.factory = factory;
        }

        /// <summary>
        /// save the package from uri to the configured save path in appsettings
        /// </summary>
        /// <param name="uri"></param>
        /// <returns></returns>
        public Task<PackageFile> DownloadAsync(PackageUri uri)
        {
            return DownloadAsync(uri, new PackageDownloadOptions
            {
                SavePath = configuration.GetHomeDirectory()
            });
        }

        public async Task<PackageFile> DownloadAsync(PackageUri uri, PackageDownloadOptions options)
        {
            var httpResult = await client.GetAsync(uri);
            var file = await DownloadFile(httpResult, options);

            return file;
        }

        private async Task<PackageFile> DownloadFile(HttpResponseMessage httpResult, PackageDownloadOptions options)
        {
            var archiveFilePath = Path.GetFullPath(Path.Combine(options.SavePath, PackageFile.FileName));

            using var resultStream = await httpResult.Content.ReadAsStreamAsync();
            using var fileStream = File.Create(archiveFilePath);

            await resultStream.CopyToAsync(fileStream);
            await resultStream.FlushAsync();

            return factory.Create(archiveFilePath);
        }
    }
}

