using System;
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
            const string environmentVariablePlatformVersion = "WEBSITE_PLATFORM_VERSION";
			const string environmentVariableSku = "WEBSITE_SKU";

			var (PlatformVersion, Sku) = (
                configuration.GetValue<string>(environmentVariablePlatformVersion),
				configuration.GetValue<string>(environmentVariableSku)
			);

			return !string.IsNullOrEmpty(PlatformVersion) && !string.IsNullOrEmpty(Sku);
        }
	}
}

