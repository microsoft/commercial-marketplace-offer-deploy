using System;
using JenkinsNET;
using Microsoft.Extensions.Options;
using Modm.Engine.Jenkins;

namespace Modm.Engine.Jenkins
{
	public class JenkinsClientFactory
	{
        private readonly JenkinsOptions options;
        private readonly ApiTokenProvider apiTokenProvider;

        public JenkinsClientFactory(IOptions<JenkinsOptions> options, ApiTokenProvider apiTokenProvider)
		{
            this.options = options.Value;
            this.apiTokenProvider = apiTokenProvider;
        }

        public async Task<IJenkinsClient> CreateAsync()
        {
            // to start making calls to Jenkins, an API Token is required. Fetch this token using the provider
            var apiToken = await apiTokenProvider.GetAsync();

            var client = new JenkinsClient
            {
                BaseUrl = options.BaseUrl,
                UserName = options.UserName,
                ApiToken = apiToken,
            };

            return client;
        }
    }
}

