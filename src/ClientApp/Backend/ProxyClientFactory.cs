// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System.Net.Http.Headers;
using ClientApp.Extensions;
using Microsoft.Extensions.DependencyInjection;

namespace ClientApp.Backend
{
    public class ProxyClientFactory
    {
        public const string BackendUrlSettingName = "BackendUrl";
        public const string ProxyClientSettingName = "ProxyClient";

        private readonly IServiceProvider provider;
        private readonly IConfiguration configuration;
        private readonly HttpClient client;
        private readonly string clientType;

        private static IProxyClient jsonFileClient;

        public ProxyClientFactory(IServiceProvider provider, IConfiguration configuration, HttpClient client)
        {
            this.provider = provider;
            this.configuration = configuration;
            this.clientType = GetClientType(configuration);
            this.client = client;
            this.client.BaseAddress = GetBaseAddress(configuration, clientType);
        }

        /// <summary>
        /// Creates a proxy client. Use the <see cref="ProxyClientSettingName"/> value in appSettings or set an environment
        /// variable to either Http or InMemory to control what is returned
        /// </summary>
        /// <param name="request">The current HttpRequest</param>
        /// <returns>and instance of a proxy client</returns>
        public IProxyClient Create(HttpRequest request)
        {
            if (ProxyClientType.IsJsonFile(clientType))
            {
                if (jsonFileClient == null)
                {
                    var logger = provider.GetRequiredService<ILogger<JsonFileProxyClient>>();
                    jsonFileClient = new JsonFileProxyClient(logger, clientType);
                }

                return jsonFileClient;
            }

            client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", request.GetJwtToken());
            return new HttpProxyClient(client, provider.GetRequiredService<ILogger<HttpProxyClient>>());
        }

        /// <summary>
        /// Gets the backend url as a base address to be used by the http client.
        /// This value MUST exist if the proxy client is Http
        /// </summary>
        /// <param name="configuration"></param>
        /// <returns></returns>
        /// <exception cref="InvalidOperationException"></exception>
        private static Uri GetBaseAddress(IConfiguration configuration, string clientType)
        {
            var value = configuration[BackendUrlSettingName]?.TrimEnd('/');

            if (string.IsNullOrEmpty(value) && ProxyClientType.IsHttp(clientType))
            {
                throw new InvalidOperationException($"Base URL of the backend is required. Setting Name is {BackendUrlSettingName}");
            }

            // trailing slash is required for a valid base address
            return new Uri(value + "/");
        }

        private static string GetClientType(IConfiguration configuration)
        {
            var value = configuration[ProxyClientSettingName] ?? "Http";
            return value;
        }
    }
}