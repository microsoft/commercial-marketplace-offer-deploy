using Modm.Azure.Model;

namespace Modm.Azure
{
    public interface IManagedIdentityService
    {
        /// <summary>
        /// Is the service endpoint reachable
        /// </summary>
        /// <returns></returns>
        /// <remarks>
        /// if the request to the IMDS is not successful, then a managed identity isn't assigned to the
        /// vm instance and we don't have an accessible identity to consume
        /// </remarks>
        Task<bool> IsAccessibleAsync(CancellationToken cancellationToken = default);
        Task<ManagedIdentityInfo> GetAsync(CancellationToken cancellationToken = default);
    }
}