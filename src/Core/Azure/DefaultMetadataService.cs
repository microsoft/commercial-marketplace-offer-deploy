using System.Net;
using System.Text.Json;
using Modm.Azure.Model;
using Polly;
using Polly.Retry;

namespace Modm.Azure
{
    /// <summary>
    /// Consuming service of the Azure Metadata Service
    /// </summary>
    /// <remarks>
    /// This will only work in an Azure VM that's running inside a vnet since hte ImdsServer endpoint is only
    /// available there.
    /// </remarks>
    /// <see cref="https://learn.microsoft.com/en-us/azure/virtual-machines/instance-metadata-service?tabs=linux#usage"/>
    public class DefaultMetadataService : IMetadataService
    {
        private const string DefaultApiVersion = "2021-02-01";
        const string DefaultServiceEndpoint = "http://169.254.169.254";
        const string InstanceEndpoint = DefaultServiceEndpoint + "/metadata/instance";

        private readonly HttpClient client;
        private readonly AsyncRetryPolicy asyncRetryPolicy;

        public DefaultMetadataService(HttpClient client)
        {
            this.client = client;
            this.asyncRetryPolicy = Policy
               .Handle<HttpRequestException>()
               .WaitAndRetryForeverAsync(retryAttempt => TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)));
        }

        public async Task<InstanceMetadata> GetAsync()
        {
            var metdata = await GetAsync<InstanceMetadata>(InstanceEndpoint, DefaultApiVersion);

            if (metdata == null)
            {
                throw new NullReferenceException("Metadata about the instance is null.");
            }

            return metdata;
        }

        public async Task<string> GetFqdnAsync()
        {
            var metadata = await GetAsync();
            var dnsLabel = $"modm{ArmFunctions.UniqueString(metadata.Compute.ResourceId)}";
            return $"{dnsLabel}.{metadata.Compute.Location}.cloudapp.azure.com";
        }

        private async Task<T> GetAsync<T>(string uri, string apiVersion, string otherParams = default)
        {
            var result = await asyncRetryPolicy.ExecuteAsync(async () => {
                var requestUri = uri + "?api-version=" + apiVersion;
                if (!string.IsNullOrEmpty(otherParams))
                {
                    requestUri += "&" + otherParams;
                }

                // IMDS requires bypassing proxies.
                HttpClient.DefaultProxy = new WebProxy();
                var request = new HttpRequestMessage(HttpMethod.Get, requestUri);
                request.Headers.Add("Metadata", "True");

                var response = await client.SendAsync(request);
                response.EnsureSuccessStatusCode();

                return await JsonSerializer.DeserializeAsync<T>(response.Content.ReadAsStream());
            });

            return result;
        }
    }
}