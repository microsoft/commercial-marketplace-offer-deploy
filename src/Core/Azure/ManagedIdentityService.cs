using System;
using System.IdentityModel.Tokens.Jwt;
using System.Net;
using System.Text.Json;
using Azure.Core;

namespace Modm.Azure
{
    /// <summary>
    /// Gets managed identity information using Imds
    /// </summary>
    /// <remarks>
    /// <see cref="https://learn.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/how-to-use-vm-token#get-a-token-using-c"/>
    /// </remarks>
	public class ManagedIdentityService
	{
		public const string TokenEndpoint = "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/";
        private readonly HttpClient client;
        private readonly InstanceMetadataService metadataService;

        public ManagedIdentityService(HttpClient client, InstanceMetadataService metadataService)
		{
            this.client = client;
            this.metadataService = metadataService;
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

