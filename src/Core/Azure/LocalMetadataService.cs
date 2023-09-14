using System;

namespace Modm.Azure
{
    public class LocalMetadataService : IMetadataService
	{
		public LocalMetadataService()
		{
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
                    Offer = "",
                    OsProfile = new OsProfile { AdminUsername = "", ComputerName = "", DisablePasswordAuthentication = false },
                    OsType = "",
                    PlacementGroupId = "",
                    Plan = new Plan { Name = "", Product = "", Publisher = "" },
                    Priority = "",
                    Provider = "",
                    Publisher = "",
                    ResourceGroupName = "",
                    ResourceId = "",
                    Sku = "",
                    StorageProfile = new StorageProfile {
                        DataDisks = Array.Empty<object>(),
                        ImageReference = new ImageReference { Id = "", Offer = "", Publisher = "", Sku = "", Version = "" },
                        OsDisk = new OsDisk { Caching = "", CreateOption = "", Name = "", OsType = "" }
                    },
                    Tags = "",
                    TagsList = new List<KeyValuePair<string, string>>(),
                    UserData = "",
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

