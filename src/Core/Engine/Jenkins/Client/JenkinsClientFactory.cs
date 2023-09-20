using System;
using Microsoft.Extensions.Options;
using Polly.Retry;

namespace Modm.Engine.Jenkins.Client
{
	class JenkinsClientFactory
	{
        private readonly JenkinsOptions options;
        private readonly System.Net.Http.HttpClient httpClient;
        private readonly ApiTokenClient apiTokenClient;
        private readonly AsyncRetryPolicy retryPolicy;

        public JenkinsClientFactory(System.Net.Http.HttpClient httpClient, ApiTokenClient apiTokenClient, IOptions<JenkinsOptions> options, AsyncRetryPolicy retryPolicy)
		{
            this.options = options.Value;
            this.httpClient = httpClient;
            this.apiTokenClient = apiTokenClient;
            this.retryPolicy = retryPolicy;
        }

        public async Task<IJenkinsClient> Create()
        {
            // to start making calls to Jenkins, an API Token is required. Fetch this token using the provider
            var apiToken = await apiTokenClient.Get();

            var client = new JenkinsClient(httpClient, options, retryPolicy)
            {
                BaseUrl = options.BaseUrl,
                UserName = options.UserName,
                ApiToken = apiToken,
            };

            return client;
        }
    }
}

