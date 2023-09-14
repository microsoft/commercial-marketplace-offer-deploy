using System;

namespace Modm.Azure
{
	public class LocalManagedIdentityService : IManagedIdentityService
	{
		public LocalManagedIdentityService()
		{
		}

        public Task<ManagedIdentityInfo> GetAsync()
        {
            return Task.FromResult(new ManagedIdentityInfo
            {
                ClientId = Guid.Empty,
                SubscriptionId = Guid.Empty,
                ObjectId = Guid.Empty,
                TenantId = Guid.Empty
            });
        }
    }
}

