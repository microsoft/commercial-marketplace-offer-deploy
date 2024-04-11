using Microsoft.Azure.Management.ResourceManager.Fluent.Core;
using Microsoft.Extensions.Configuration;

namespace Modm.Configuration
{
    /// <summary>
    /// uses the environment variables provided by app service to determine the app configuration resource
    /// <see cref="https://learn.microsoft.com/en-us/azure/app-service/reference-app-settings?tabs=kudu%2Cdotnet"/>
    ///
    /// This depends on <see cref="IConfiguration"/> which must be loaded with environment variables using
    /// builder.Configuration.AddEnvironmentVariables();
    /// </summary>
    public class AppServiceAppConfigurationResourceProvider : IAppConfigurationResourceProvider
	{
        public const string EnvironmentVariable_OwnerName = "WEBSITE_OWNER_NAME";
        public const string EnvironmentVariable_ResourceGroupName = "WEBSITE_RESOURCE_GROUP";

        private readonly string subscriptionId;
        private readonly string resourceGroupName;

        public AppServiceAppConfigurationResourceProvider(IConfiguration configuration)
		{
            var websiteOwnerName = configuration[EnvironmentVariable_OwnerName];

            this.subscriptionId = websiteOwnerName[..websiteOwnerName.IndexOf('+')];
            this.resourceGroupName = configuration[EnvironmentVariable_ResourceGroupName];
		}

        public AppConfigurationResource Get()
        {
            var resourceGroupId = ResourceId.FromString($"/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}");
            return new AppConfigurationResource(resourceGroupId);
        }
    }
}