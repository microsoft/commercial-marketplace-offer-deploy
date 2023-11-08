using Azure.Identity;
using Azure.ResourceManager;
using Microsoft.Extensions.Azure;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Modm.Azure;
using Modm.Extensions;
using Microsoft.Azure.Management.ResourceManager;
using Azure.Core;

namespace Modm.Configuration
{
    /// <summary>
    /// used for bootstrapping the appconfiguration by providing the appconfigstore endpoint
    /// using the Imds
    /// </summary>
    public class AppConfigurationRegistration
	{
        private readonly IMetadataService metadataService;
        private readonly ArmClient client;
        private readonly IConfigurationBuilder builder;

        public AppConfigurationRegistration(IMetadataService metadataService, ArmClient client, IConfigurationBuilder builder)
		{
            this.metadataService = metadataService;
            this.client = client;
            this.builder = builder;
        }

        public void AddAppConfigurationIfExists()
        {
            var metadata = metadataService.GetAsync().GetAwaiter().GetResult();
            var resource = new AppConfigurationResource(metadata.ResourceGroupId);

            var resourceOperations = client.GetGenericResource(resource.Identifier);

            try
            {
                var azureResource = resourceOperations.Get();

                if (azureResource != null && !azureResource.GetRawResponse().IsError)
                {
                    builder.AddAzureAppConfiguration(options => options.Connect(resource.Uri, new DefaultAzureCredential()));
                }
            }
            catch (Exception)
            {

            }
        }
    }
}

