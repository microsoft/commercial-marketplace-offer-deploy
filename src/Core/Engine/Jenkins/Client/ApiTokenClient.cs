using System.Net.Http.Headers;
using System.Text.Json;
using Microsoft.Extensions.Options;
using Modm.Engine.Jenkins.Model;
using Modm.Extensions;

namespace Modm.Engine.Jenkins.Client
{
    /// <summary>
    /// Will perform the job of getting a Jenkins API Token
    /// </summary>
    /// <remarks>
    /// In order to do anything with Jenkins API, an API token is required.
    /// To get an API token, the username:password is used to authenticate.
    /// Get a crumb, use the crumb to create an API Token, then return the token value
    /// </remarks>
	public class ApiTokenClient
	{
        private readonly HttpClient client;
        private readonly JenkinsOptions options;

        public ApiTokenClient(HttpClient client, IOptions<JenkinsOptions> options)
		{
            this.client = client;
            this.options = options.Value;
        }

        /// <summary>
        /// Gets the API token using the configured <see cref="JenkinsOptions"/>
        /// </summary>
        /// <param name="options"></param>
        /// <returns></returns>
        /// <exception cref="JenkinsCrumbRequestException">If the crumb response is null</exception>
        public async Task<string> Get()
        {
            var crumbResponse = await GetCrumb() ?? throw new JenkinsCrumbRequestException("Crumb response null");
            var response = await GenerateApiToken(crumbResponse);

            if (response != null && response.Status == "ok")
            {
                // TODO: start logging things from HTTP client responses centrally to capture flow
                return response.Data.Value;
            }

            throw new InvalidOperationException("Generate API Token response is null or status was not 'ok'");
        }

        private async Task<GenerateApiTokenResponse> GenerateApiToken(GetCrumbResponse crumbResponse)
        {
            using var request = GetHttpRequest(HttpMethod.Post, "me/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken");
            request.Headers.Add(crumbResponse.RequestField, crumbResponse.Crumb);

            var response = await client.SendAsync(request);
            response.EnsureSuccessStatusCode();

            return await Deserialize<GenerateApiTokenResponse>(response);
        }

        private async Task<GetCrumbResponse> GetCrumb()
        {
            using var request = GetHttpRequest(HttpMethod.Get, "crumbIssuer/api/json");
            var response = await client.SendAsync(request);

            try
            {
                response.EnsureSuccessStatusCode();
            }
            catch (Exception ex)
            {
                throw new JenkinsCrumbRequestException("Received a non success status coding while assempting to get crumb", ex);
            }

            return await Deserialize<GetCrumbResponse>(response);
        }

        private HttpRequestMessage GetHttpRequest(HttpMethod method, string relativeUri)
        {
            var requestUri = new Uri(options.BaseUrl).Append(relativeUri).AbsoluteUri;
            var request = new HttpRequestMessage(method, requestUri);

            request.Headers.Authorization = GetAuthenticationHeader(options);

            return request;
        }

        private static AuthenticationHeaderValue GetAuthenticationHeader(JenkinsOptions options)
        {
            var plainTextBytes = System.Text.Encoding.UTF8.GetBytes($"{options.UserName}:{options.Password}");
            var credential = Convert.ToBase64String(plainTextBytes);

            return new AuthenticationHeaderValue("Basic", credential);
        }

        private static async Task<T> Deserialize<T>(HttpResponseMessage response)
        {
            return await JsonSerializer.DeserializeAsync<T>(response.Content.ReadAsStream());
        }

    }

    
}

