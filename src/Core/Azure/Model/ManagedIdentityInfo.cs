using System;
using Azure.Core;

namespace Modm.Azure.Model
{
	public class ManagedIdentityInfo
	{
		public Guid ClientId { get; set; }
		public Guid ObjectId { get; set; }
		public Guid TenantId { get; set; }
		public Guid SubscriptionId { get; set; }
	}
}

