using System;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;

namespace Modm.Artifacts
{
	public class ArtifactsFileFactory
	{
        private readonly IServiceProvider provider;

        public ArtifactsFileFactory(IServiceProvider provider)
		{
            this.provider = provider;
        }

		public ArtifactsFile Create(string filePath)
		{
			return new ArtifactsFile(
				filePath,
				provider.GetRequiredService<ILogger<ArtifactsFile>>()
			);
		}
	}
}

