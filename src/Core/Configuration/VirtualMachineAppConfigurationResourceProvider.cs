using Modm.Azure;

namespace Modm.Configuration
{
    public class VirtualMachineAppConfigurationResourceProvider : IAppConfigurationResourceProvider
	{
        private readonly IMetadataService metadataService;

        public VirtualMachineAppConfigurationResourceProvider(IMetadataService metadataService)
		{
            this.metadataService = metadataService;
        }

        public AppConfigurationResource Get()
        {
            var metadata = metadataService.GetAsync().GetAwaiter().GetResult();
            var resource = new AppConfigurationResource(metadata.ResourceGroupId);

            return resource;
        }
    }
}