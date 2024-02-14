// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using Azure.ResourceManager.Resources;

namespace Modm.Azure
{
    public interface IAzureResourceManagerClient
    {
        Task<List<GenericResource>> GetResourcesToDeleteAsync(string resourceGroupName, string phase);
        Task<ResourceGroupResource> GetResourceGroupResourceAsync(string resourceGroupName);
        Task<bool> TryDeleteResourceAsync(GenericResource resource);
    }
}

