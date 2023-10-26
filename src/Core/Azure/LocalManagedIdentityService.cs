using System;
using Modm.Azure.Model;

namespace Modm.Azure
{
	public class LocalManagedIdentityService : IManagedIdentityService
	{
		public LocalManagedIdentityService()
		{
		}

        public Task<ManagedIdentityInfo> GetAsync(CancellationToken cancellationToken = default)
        {
            return Task.FromResult(new ManagedIdentityInfo
            {
                ClientId = Guid.Empty,
                SubscriptionId = Guid.Empty,
                ObjectId = Guid.Empty,
                TenantId = Guid.Empty
            });
        }

        public Task<bool> IsAccessibleAsync(CancellationToken cancellationToken = default)
        {
            return Task.FromResult(true);
        }
    }
}

