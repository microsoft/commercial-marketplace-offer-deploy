namespace Modm.Azure
{
    public interface IManagedIdentityService
    {
        Task<ManagedIdentityInfo> GetAsync();
    }
}