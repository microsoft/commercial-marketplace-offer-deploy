using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;

namespace Modm.Engine.Jenkins.Client
{
    class JenkinsClientFactory
	{
        private readonly JenkinsOptions options;
        private readonly HttpClient httpClient;
        private readonly ApiTokenClient apiTokenClient;
        private readonly ServiceProvider serviceProvider;

        public JenkinsClientFactory(HttpClient client, ApiTokenClient apiTokenClient, ServiceProvider serviceProvider, IOptions<JenkinsOptions> options)
		{
            this.options = options.Value;
            this.httpClient = client;
            this.apiTokenClient = apiTokenClient;
            this.serviceProvider = serviceProvider;
        }

        public async Task<IJenkinsClient> Create()
        {
            // to start making calls to Jenkins, an API Token is required. Fetch this token using the provider
            var apiToken = await apiTokenClient.Get();
            var jenkinsNetClient = new JenkinsNET.JenkinsClient(options.BaseUrl)
            {
                UserName = options.UserName,
                ApiToken = apiToken
            };

            var logger = serviceProvider.GetRequiredService<ILogger<JenkinsClient>>();
            var client = new JenkinsClient(httpClient, jenkinsNetClient, logger, options);

            return client;
        }
    }
}

