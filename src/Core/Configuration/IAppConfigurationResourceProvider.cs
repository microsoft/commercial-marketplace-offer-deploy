// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿namespace Modm.Configuration
{
    public interface IAppConfigurationResourceProvider
	{
        /// <summary>
        /// Gets the <see cref="AppConfigurationResource"/> for the resource group for the environment this is executing. For example,
        /// If this is running in a VM, it will use the VM's metadata service. If app service, the environment variables
        /// </summary>
        /// <returns>The resource id of the resource group</returns>
        AppConfigurationResource Get();
	}
}

