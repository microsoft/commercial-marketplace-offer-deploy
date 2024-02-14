// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System.IO;
using System.Text.Json;
using Microsoft.AspNetCore.Mvc;
using Modm.Deployments;
using Modm.Engine;

namespace ClientApp.Backend
{
    /// <summary>
    /// For Development purposes and frontend debugging purposes only
    /// </summary>
    public class JsonFileProxyClient : IProxyClient
    {
        public const string DefaultFileName = "proxyclient.json";

        private readonly ILogger<JsonFileProxyClient> logger;
        private FileSystemWatcher watcher;
        private readonly string filePath;
        private JsonDocument json;

        public JsonFileProxyClient(ILogger<JsonFileProxyClient> logger, string clientType)
        {
            this.logger = logger;
            this.filePath = GetFilePath(clientType);
            this.json = JsonDocument.Parse(File.ReadAllText(filePath));


            this.watcher = new FileSystemWatcher
            {
                Path = Path.GetDirectoryName(filePath),
                NotifyFilter = NotifyFilters.LastWrite,
                Filter = Path.GetFileName(filePath),
                EnableRaisingEvents = true
            };
            watcher.Changed += new FileSystemEventHandler(OnFileChanged);
        }

        private void OnFileChanged(object sender, FileSystemEventArgs e)
        {
            json = JsonDocument.Parse(File.ReadAllText(e.FullPath));
        }

        public Task<IActionResult> GetAsync<T>(string relativeUri)
        {
            logger.LogInformation("Backend request: {relativeUri}", relativeUri);
            object value = new { };

            var routes = json.RootElement.GetProperty("routes");

            try
            {
                value = routes.GetProperty(relativeUri);
            }
            catch (KeyNotFoundException ex)
            {
                logger.LogWarning(ex, "Route {relativeUri} doesn't exist in {filePath}", relativeUri, filePath);
            }

            return Task.FromResult<IActionResult>(new OkObjectResult(value));
        }

        public Task<IActionResult> PostAsync(string relativeUri, HttpContent content = null)
        {
            logger.LogInformation("Backend request: {relativeUri}", relativeUri);
            return Task.FromResult<IActionResult>(new OkObjectResult(new { }));
        }

        private static string GetFilePath(string clientType)
        {
            if (clientType.Contains('|'))
            {
                return Path.Combine(AppDomain.CurrentDomain.BaseDirectory, clientType.Split('|')[1]);
            }

            return Path.Combine(AppDomain.CurrentDomain.BaseDirectory, DefaultFileName);
        }
    }
}

