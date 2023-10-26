using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;

namespace Modm.Jenkins.Client
{
    public class JenkinsClientFactory
	{
        private readonly JenkinsOptions options;
        private readonly HttpClient httpClient;
        private readonly ApiTokenClient apiTokenClient;
        private readonly IServiceProvider serviceProvider;

        private readonly Lazy<Task<string>> apiToken;

        /// <summary>
        /// Constructor without params only to support testing
        /// </summary>
        public JenkinsClientFactory()
        {
        }

        public JenkinsClientFactory(HttpClient client, ApiTokenClient apiTokenClient, IServiceProvider serviceProvider, IOptions<JenkinsOptions> options)
		{
            this.options = options.Value;
            this.httpClient = client;
            this.apiTokenClient = apiTokenClient;
            this.serviceProvider = serviceProvider;
            this.apiToken = new(apiTokenClient.Get, LazyThreadSafetyMode.PublicationOnly);
        }

        public virtual async Task<IJenkinsClient> Create()
        {
            var apiToken = await GetApiToken();
            var jenkinsNetClient = new JenkinsNET.JenkinsClient(options.BaseUrl)
            {
                UserName = options.UserName,
                ApiToken = apiToken
            };

            // add the api token to the options 

            var logger = serviceProvider.GetRequiredService<ILogger<JenkinsClient>>();
            var client = new JenkinsClient(httpClient, jenkinsNetClient, logger, options);

            return client;
        }

        /// <summary>
        /// to start making calls to Jenkins, an API Token is required. Fetch this token using the provider
        /// </summary>
        /// <returns></returns>
        public async ValueTask<string> GetApiToken()
        {
            var value = await apiToken.Value.ConfigureAwait(false);
            return value;
        }
    }
}

