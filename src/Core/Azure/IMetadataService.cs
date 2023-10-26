using Modm.Azure.Model;

namespace Modm.Azure
{
    public interface IMetadataService
    {
        Task<InstanceMetadata> GetAsync();
        Task<string> GetFqdnAsync();
    }
}