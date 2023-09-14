namespace Modm.Azure
{
    public interface IManagedIdentityService
    {
        /// <summary>
        /// Is the service endpoint reachable
        /// </summary>
        /// <returns></returns>
        Task<bool> IsAccessibleAsync();
        Task<ManagedIdentityInfo> GetAsync();
    }
}