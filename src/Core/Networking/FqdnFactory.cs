using System;
using Modm.Azure;

namespace Modm.Networking
{
	public static class FqdnFactory
	{
        /// <summary>
        /// Returns the FQDN of the machine using the machine name and the location
        /// </summary>
        /// <remarks>
        /// We MUST have this value set on the VM's NIC dnsLabel in order to launch the containers, required by Caddy
        /// Format: {uniqueString(vm.resourceId)}.{location}.cloudapp.azure.com
        /// </remarks>
        /// <returns></returns>
        public static string Create(string resourceId, string location)
		{
			var dnsLabel = ArmFunctions.UniqueString(resourceId);
            return $"{dnsLabel}.{location}.cloudapp.azure.com";
        }
	}
}

