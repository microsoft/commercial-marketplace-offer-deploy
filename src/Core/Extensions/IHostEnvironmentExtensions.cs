using System;
using Microsoft.Extensions.Hosting;

namespace Modm.Extensions
{
	public static class IHostEnvironmentExtensions
	{
        /// <summary>
        /// We need this otherwise during the packer build, the servicehost will try and
        /// connect to a non-existent app configuration resource
        /// </summary>
        /// <remarks>
        ///	This works through the packer file, setting the DOTNET_ENVIRONMENT in the packer provisioner
        /// </remarks>
        /// <param name="environment"></param>
        /// <returns></returns>
        public static bool IsPacker(this IHostEnvironment environment)
		{
			return environment.EnvironmentName == "Packer";
		}
	}
}

