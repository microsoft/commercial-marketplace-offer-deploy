using System;
using JenkinsNET;
using Microsoft.Extensions.Options;

namespace Modm.Engine
{
	public class JenkinsClientFactory
	{
        private readonly JenkinsOptions options;

        public JenkinsClientFactory(IOptions<JenkinsOptions> options)
		{
            this.options = options.Value;
        }

        public IJenkinsClient Create()
        {
            var client = new JenkinsClient
            {
                BaseUrl = options.BaseUrl,
                UserName = options.UserName,
                ApiToken = options.ApiToken,
            };

            return client;
        }
    }
}

