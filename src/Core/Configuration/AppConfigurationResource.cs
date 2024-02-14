// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using Modm.Azure;
using Microsoft.Azure.Management.ResourceManager.Fluent.Core;
using Azure.Core;

namespace Modm.Configuration
{
    /// <summary>
    /// Recreation of what's in the ./templates/mainTemplate/variables
    /// This is required in order to ensure that we can reliably get the app configuration
    /// store name without having to pass it to the VM through userData
    /// 
    /// [concat('modmconfig-', substring(uniqueString(resourceGroup().id), 0, 8))]
    /// </summary>
    public struct AppConfigurationResource
	{
        private const string prefix = "modmconfig-";
        private const int uniqueSuffixLength = 8;

        private readonly string suffix;

        public string Name { get; }

        public Uri Uri { get; }

        public ResourceIdentifier Identifier { get; }

		public AppConfigurationResource(ResourceId resourceGroupId)
		{
            this.suffix = ArmFunctions.UniqueString(resourceGroupId.Id).Substring(0, uniqueSuffixLength);
            this.Name = string.Concat(prefix, suffix);
            this.Uri = new Uri($"https://{Name}.azconfig.io");
            this.Identifier = ResourceIdentifier.Parse($"{resourceGroupId.Id}/providers/Microsoft.AppConfiguration/configurationStores/{Name}");
        }
	}
}