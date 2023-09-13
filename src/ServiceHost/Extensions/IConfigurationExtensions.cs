using System;
using Modm.Configuration;

namespace Modm.ServiceHost.Extensions
{
	public static class IConfigurationExtensions
	{
		public static string GetHomeDirectory(this IConfiguration configuration)
		{
			var value = configuration.GetValue<string>(EnvironmentVariable.Names.HomeDirectory);
			return value ?? string.Empty;
		}
	}
}

