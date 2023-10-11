using System;
using Microsoft.Extensions.Configuration;
using Modm.Azure.Model;

namespace Modm.Azure
{
    public class LocalMetadataService : IMetadataService
	{
        private readonly IConfiguration configuration;

        public LocalMetadataService(IConfiguration configuration)
		{
            this.configuration = configuration;
        }

        public Task<InstanceMetadata> GetAsync()
        {
            return Task.FromResult(new InstanceMetadata {
                Compute = new Compute
                {
                    AzEnvironment = "",
                    EvictionPolicy = "",
                    CustomData = "",
                    LicenseType = "",
                    Location = "",
                    Name = "",
                    Offer = "Local Offer",
                    OsProfile = new OsProfile { AdminUsername = "", ComputerName = "", DisablePasswordAuthentication = false },
                    OsType = "",
                    PlacementGroupId = "",
                    Plan = new Plan { Name = "", Product = "", Publisher = "" },
                    Priority = "",
                    Provider = "",
                    Publisher = "",
                    ResourceGroupName = configuration.GetSection("Azure").GetValue<string>("DefaultResourceGroupName"),
                    SubscriptionId = string.IsNullOrEmpty(configuration.GetSection("Azure").GetValue<string>("DefaultSubscriptionId"))
                        ? Guid.Empty
                        : Guid.Parse(configuration.GetSection("Azure").GetValue<string>("DefaultSubscriptionId")),
                    ResourceId = "",
                    Sku = "",
                    StorageProfile = new StorageProfile {
                        DataDisks = Array.Empty<object>(),
                        ImageReference = new ImageReference { Id = "", Offer = "", Publisher = "", Sku = "", Version = "" },
                        OsDisk = new OsDisk { Caching = "", CreateOption = "", Name = "", OsType = "" }
                    },
                    Tags = "",
                    TagsList = new List<KeyValuePair<string, string>>(),
                    UserData = new UserData { ArtifactsUri = Environment.GetEnvironmentVariable("ARTIFACTS_URL") ?? "", ArtifactsHash = Environment.GetEnvironmentVariable("ARTIFACTS_HASH") ?? "" }.ToBase64Json(),
                    Version = "",
                    VmScaleSetName = "",
                    VmSize = "",
                    Zone = ""
                },
                Network = new Network() });
        }

        public Task<string> GetFqdnAsync()
        {
            return Task.FromResult("localhost");
        }
    }
}

