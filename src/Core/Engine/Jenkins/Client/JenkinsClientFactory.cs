using Microsoft.Extensions.Options;
using Modm.Http;

namespace Modm.Engine.Jenkins.Client
{
    class JenkinsClientFactory
	{
        private readonly JenkinsOptions options;
        private readonly HttpClient httpClient;
        private readonly ApiTokenClient apiTokenClient;

        public JenkinsClientFactory(IHttpClientFactory clientFactory, ApiTokenClient apiTokenClient, IOptions<JenkinsOptions> options)
		{
            this.options = options.Value;
            this.httpClient = clientFactory.CreateClient(HttpConstants.DefaultHttpClientName);
            this.apiTokenClient = apiTokenClient;
        }

        public async Task<IJenkinsClient> Create()
        {
            // to start making calls to Jenkins, an API Token is required. Fetch this token using the provider
            var apiToken = await apiTokenClient.Get();

            var client = new JenkinsClient(httpClient, options)
            {
                BaseUrl = options.BaseUrl,
                UserName = options.UserName,
                ApiToken = apiToken,
            };

            return client;
        }
    }
}

