// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System.IdentityModel.Tokens.Jwt;
using System.Net;
using System.Text.Json;
using Microsoft.Extensions.Logging;
using Modm.Azure.Model;

namespace Modm.Azure
{
    /// <summary>
    /// Gets managed identity information using Imds
    /// </summary>
    /// <remarks>
    /// <see cref="https://learn.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/how-to-use-vm-token#get-a-token-using-c"/>
    /// </remarks>
	public class DefaultManagedIdentityService : IManagedIdentityService
    {
        public const string TokenEndpoint = "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/";

        /// <summary>
        /// empty web proxy will bypass proxies which is required by IMDS
        /// </summary>
        static readonly WebProxy ByPassWebProxy = new();

        private readonly HttpClient client;
        private readonly IMetadataService metadataService;
        private readonly ILogger<DefaultManagedIdentityService> logger;

        public DefaultManagedIdentityService(HttpClient client, IMetadataService metadataService, ILogger<DefaultManagedIdentityService> logger)
        {
            this.client = client;
            this.metadataService = metadataService;
            this.logger = logger;
        }

        public async Task<ManagedIdentityInfo> GetAsync(CancellationToken cancellationToken = default)
        {
            var metadata = await metadataService.GetAsync();
            var token = await GetTokenAsync(cancellationToken);

            if (token == null)
            {
                throw new InvalidOperationException("Token was null while fetching managed identity information.");
            }

            var accessToken = new JwtSecurityToken(token.AccessToken);

            return new ManagedIdentityInfo
            {
                ClientId = token.ClientId,
                SubscriptionId = metadata.Compute.SubscriptionId,
                TenantId = Guid.Parse(accessToken.Claims.First(c => c.Type == "oid").Value),
                ObjectId = Guid.Parse(accessToken.Claims.First(c => c.Type == "tid").Value)
            };
        }

        public async Task<bool> IsAccessibleAsync(CancellationToken cancellationToken = default)
        {
            var request = CreateRequest();
            var response = await client.SendAsync(request, cancellationToken);

            if (response.StatusCode == HttpStatusCode.BadRequest || !response.IsSuccessStatusCode)
            {
                var result = await JsonSerializer.DeserializeAsync<AcquireTokenErrorResponse>(
                    response.Content.ReadAsStream(cancellationToken),
                    cancellationToken: cancellationToken
                    );
                logger.LogTrace("Error received while checking IMDS access token endpoint: {message}.", result?.ErrorDescription);
                return false;
            }

            return true;
        }

        private async Task<AcquireTokenResponse> GetTokenAsync(CancellationToken cancellationToken)
        {
            var request = CreateRequest();
            var response = await client.SendAsync(request, cancellationToken);

            try
            {
                response.EnsureSuccessStatusCode();

                var result = await JsonSerializer.DeserializeAsync<AcquireTokenResponse>(
                    response.Content.ReadAsStream(cancellationToken),
                    cancellationToken: cancellationToken);
                return result;

            }
            catch (Exception ex)
            {
                logger.LogError(ex, "Failed to receive access token response from metadata service endpoint for authentication token");
            }

            return null;
        }

        private static HttpRequestMessage CreateRequest()
        {
            HttpClient.DefaultProxy = ByPassWebProxy;
            var request = new HttpRequestMessage(HttpMethod.Get, TokenEndpoint);
            request.Headers.Add("Metadata", "True");

            return request;
        }
    }
}