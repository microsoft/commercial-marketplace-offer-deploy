// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using Microsoft.Extensions.Configuration;
using Modm.Configuration;

namespace Modm.Extensions
{
	public static class IConfigurationExtensions
	{
		/// <summary>
		/// Get home directory of MODM (resolves to value in $MODM_HOME)
		/// </summary>
		/// <param name="configuration"></param>
		/// <returns></returns>
		public static string GetHomeDirectory(this IConfiguration configuration)
		{
			var value = configuration.GetValue<string>(EnvironmentVariable.Names.HomeDirectory);
			return value ?? string.Empty;
		}

		/// <summary>
		/// Whether this is an app service environment configuration
		/// </summary>
		/// <param name="configuration"></param>
		/// <returns></returns>
		public static bool IsAppServiceEnvironment(this IConfiguration configuration)
		{
            // determine based on whether the env is present, which is provided by app service
            // see: https://learn.microsoft.com/en-us/azure/app-service/reference-app-settings?tabs=kudu%2Cdotnet
			var (RunFromPackage, Sku) = (
                configuration.GetValue<string>("WEBSITE_RUN_FROM_PACKAGE"),
				configuration.GetValue<string>("WEBSITE_SKU")
			);

			return !string.IsNullOrEmpty(RunFromPackage) && !string.IsNullOrEmpty(Sku);
        }
	}
}

