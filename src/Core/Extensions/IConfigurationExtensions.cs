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
	}
}

