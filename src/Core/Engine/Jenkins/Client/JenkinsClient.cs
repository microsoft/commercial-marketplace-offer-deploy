using System;
using System.Net.Http.Headers;
using System.Text.Json;
using Microsoft.Extensions.Options;
using Polly;
using Polly.Retry;
using Modm.Engine.Jenkins.Model;
using Modm.Extensions;

namespace Modm.Engine.Jenkins.Client
{
    class JenkinsClient : JenkinsNET.JenkinsClient, IJenkinsClient
	{
        const string JenkinsVersionHeaderName = "X-Jenkins";

        private readonly System.Net.Http.HttpClient client;
        private readonly JenkinsOptions options;
        private readonly AsyncRetryPolicy retryPolicy;

        public JenkinsClient(System.Net.Http.HttpClient client, JenkinsOptions options, AsyncRetryPolicy retryPolicy)
        {
            this.client = client;
            this.options = options;
            this.retryPolicy = retryPolicy;
        }

        public async Task<JenkinsInfo> GetInfo()
        {
            var version = "";
            var hudson = await Send<Hudson>(HttpMethod.Get, "api/json", response =>
            {
                version = response.Headers.GetValues(JenkinsVersionHeaderName).FirstOrDefault();
                response.EnsureSuccessStatusCode();
            });

            return new JenkinsInfo { Hudson = hudson, Version = version };
        }

        public async Task<MasterComputer> GetBuiltInNode()
        {
            return await Send<MasterComputer>(HttpMethod.Get, "computer/(built-in)/api/json", response => response.EnsureSuccessStatusCode());
        }

        /// <summary>
        /// Helper method that sends an HTTP request to Jenkins using the method and relative Uri with an optional http response message handler,
        /// then deserializes the JSON response
        /// </summary>
        /// <typeparam name="T">The type of response to deserialize</typeparam>
        /// <param name="method"></param>
        /// <param name="relativeUri"></param>
        /// <param name="handler"></param>
        /// <returns></returns>
        private async Task<T> Send<T>(HttpMethod method, string relativeUri, Action<HttpResponseMessage> handler = null)
        {
            return await retryPolicy.ExecuteAsync(async () =>
            {
                using var request = GetHttpRequest(method, relativeUri);
                var response = await client.SendAsync(request);

                if (handler != null)
                {
                    handler(response);
                }

                return await Deserialize<T>(response);
            });
        }

        private HttpRequestMessage GetHttpRequest(HttpMethod method, string relativeUri)
        {
            var requestUri = new Uri(options.BaseUrl).Append(relativeUri).AbsoluteUri;
            var request = new HttpRequestMessage(method, requestUri);

            request.Headers.Authorization = GetAuthenticationHeader(options);

            return request;
        }

        /// <summary>
        /// gets an authentication header using the username and api token
        /// </summary>
        /// <param name="options"></param>
        /// <returns></returns>
        private static AuthenticationHeaderValue GetAuthenticationHeader(JenkinsOptions options)
        {
            var plainTextBytes = System.Text.Encoding.UTF8.GetBytes($"{options.UserName}:{options.ApiToken}");
            var credential = Convert.ToBase64String(plainTextBytes);

            return new AuthenticationHeaderValue("Basic", credential);
        }

        /// <summary>
        /// deserializes the response as JSON
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="response"></param>
        /// <returns></returns>
        private static async Task<T> Deserialize<T>(HttpResponseMessage response)
        {
            return await JsonSerializer.DeserializeAsync<T>(response.Content.ReadAsStream());
        }
    }
}

