using Azure.Core;
using Azure.Identity;
using System.Net.Http;
using System.Net.Http.Headers;
using System.IO;
using System.Threading.Tasks;

namespace Modm.Azure
{
	public class KuduApi
    {
        private readonly HttpClient client;
        private readonly string functionAppName;
        private readonly DefaultAzureCredential credential;

        public KuduApi(string functionAppName, HttpClient client)
        {
            this.functionAppName = functionAppName;
            this.client = client;
            this.credential = new DefaultAzureCredential();
        }

        public async Task DeployZipAsync(string zipFilePath)
        {
            var baseUrl = $"https://{this.functionAppName}.scm.azurewebsites.net/api/zipdeploy/";

            using var request = new HttpRequestMessage(HttpMethod.Post, baseUrl);

            var token = await this.credential.GetTokenAsync(new TokenRequestContext(new[] { "https://management.azure.com/.default" }));
            request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", token.Token);

            using var fileStream = new FileStream(zipFilePath, FileMode.Open);
            using var streamContent = new StreamContent(fileStream);


            request.Content = streamContent;

            var response = await this.client.SendAsync(request);
            response.EnsureSuccessStatusCode();
        }
    }
}

