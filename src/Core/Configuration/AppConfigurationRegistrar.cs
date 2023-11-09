using Azure.Identity;
using Azure.ResourceManager;
using Microsoft.Extensions.Configuration;
using Modm.Azure;
using Microsoft.Extensions.Logging;
using Azure.Core;
using Azure.ResourceManager.AppConfiguration;
using Azure;
using Azure.ResourceManager.Resources;
using Azure.ResourceManager.AppConfiguration.Models;

namespace Modm.Configuration
{
    /// <summary>
    /// used for bootstrapping the appconfiguration by providing the appconfigstore endpoint
    /// using the Imds
    /// </summary>
    public class AppConfigurationRegistrar
	{
        private readonly IMetadataService metadataService;
        private readonly ArmClient client;
        private readonly IConfigurationBuilder builder;
        private readonly ILogger<AppConfigurationRegistrar> logger;

        public AppConfigurationRegistrar(IMetadataService metadataService, ArmClient client, IConfigurationBuilder builder, ILogger<AppConfigurationRegistrar> logger)
		{
            this.metadataService = metadataService;
            this.client = client;
            this.builder = builder;
            this.logger = logger;
        }

        public void AddAppConfigurationIfExists()
        {
            var metadata = metadataService.GetAsync().GetAwaiter().GetResult();
            var resource = new AppConfigurationResource(metadata.ResourceGroupId);

            using (logger.BeginScope(new Dictionary<string, object>
            {
                ["resourceGroupId"] = metadata.ResourceGroupId.Id,
                ["appConfigurationResourceId"] = resource.Identifier.ToString()
            }))
            {
                try
                {
                    var appConfigurationStoreResource = client.GetAppConfigurationStoreResource(resource.Identifier);

                    if (Exists(appConfigurationStoreResource, out AppConfigurationStoreResource appConfigurationStore))
                    {
                        var connectionString = GetConnectionString(appConfigurationStore);

                        logger.LogInformation("Registering with configuration using connection string.");
                        builder.AddAzureAppConfiguration(options => options.Connect(connectionString));

                        return;
                    }

                    logger.LogInformation("App Configuration does NOT exist, skipping configuration registration.");

                }
                catch (Exception e)
                {
                    logger.LogError(e, "Error trying to register App Configuration resource with configuration.");
                }
            }
        }

        private bool Exists(AppConfigurationStoreResource appConfigurationStoreResource, out AppConfigurationStoreResource resource)
        {
            logger.LogInformation("Checking if appConfiguration resource exists.");

            try
            {
                var response = appConfigurationStoreResource.Get();
                resource = response.Value;
            }
            catch (RequestFailedException ex) when (ex.Status == 404)
            {
                logger.LogWarning($"App Configuration {appConfigurationStoreResource.Id.Name} does not exist.");
                resource = null;
                return false;
            }

            logger.LogInformation("App Configuration exists, registering with configuration.");
            return true;
        }

        private string GetConnectionString(AppConfigurationStoreResource appConfigurationStore)
        {
            logger.LogDebug("Attempting to fetch app configuration store keys");

            var pageableKeys = appConfigurationStore.GetKeys();
            foreach (var key in pageableKeys)
            {
                if (!key.IsReadOnly.GetValueOrDefault(true))
                {
                    logger.LogInformation("ConnectionString found.");
                    return key.ConnectionString;
                }
            }

            throw new InvalidOperationException($"Failed to get an access key with a connection string for {appConfigurationStore.Id.Name}");
        }
    }
}