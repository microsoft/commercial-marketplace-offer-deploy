using System;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;

namespace Modm.Packaging
{
	public class PackageFileFactory
	{
        private readonly IServiceProvider provider;

        public PackageFileFactory(IServiceProvider provider)
		{
            this.provider = provider;
        }

		public PackageFile Create(string filePath)
		{
			return new PackageFile(
				filePath,
				provider.GetRequiredService<ILogger<PackageFile>>()
			);
		}
	}
}

