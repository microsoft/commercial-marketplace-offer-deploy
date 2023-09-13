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

            });
        }
    }
}

