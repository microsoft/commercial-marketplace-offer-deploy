using System;
using System.IdentityModel.Tokens.Jwt;
using System.Net;
using System.Text.Json;
using Azure.Core;
using Microsoft.Extensions.Logging;

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
        private readonly HttpClient client;
        private readonly IMetadataService metadataService;
        private readonly ILogger<DefaultManagedIdentityService> logger;

        public DefaultManagedIdentityService(HttpClient client, IMetadataService metadataService, ILogger<DefaultManagedIdentityService> logger)
        {
            this.client = client;
            this.metadataService = metadataService;
            this.logger = logger;
        }

        public async Task<ManagedIdentityInfo> GetAsync()
        {
            var metadata = await metadataService.GetAsync();
            var token = await GetTokenAsync();

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

        public async Task<bool> IsAccessibleAsync()
        {
            try
            {
                var result = await GetTokenAsync();
                return (result != null && !string.IsNullOrEmpty(result.AccessToken));
            }
            catch (Exception ex)
            {
                logger.LogTrace(ex, "Managed Identity metdata service endpoint unreachable");
            }

            return false;
        }

        private async Task<TokenResponse?> GetTokenAsync()
        {
            // IMDS requires bypassing proxies.
            HttpClient.DefaultProxy = new WebProxy();

            var request = new HttpRequestMessage(HttpMethod.Get, TokenEndpoint);
            request.Headers.Add("Metadata", "True");

            var response = await client.SendAsync(request);
            response.EnsureSuccessStatusCode();
            var result = await JsonSerializer.DeserializeAsync<TokenResponse>(response.Content.ReadAsStream());


            return result;
        }
    }
}

