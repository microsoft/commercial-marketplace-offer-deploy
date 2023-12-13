using System;
using Azure.Core;

namespace Modm.Tests.Utils
{
    public partial class Test
    {
        public static ResourceIdentifier AzureResourceIdentifier(string resourceGroupName, ResourceType resourceType)
        {
            var input = string.Concat(
                "/subscriptions/", Guid.NewGuid().ToString(),
                "/resourceGroups/", resourceGroupName,
                "/providers/", resourceType.ToString(),
                "/", RandomString(5));

            return ResourceIdentifier.Parse(input);
        }
    }
}

