using System;
using Microsoft.Extensions.DependencyInjection;

namespace Modm.Deployments
{
	public class DeploymentFileFactory : IDeploymentFileFactory
	{
        private readonly IServiceProvider serviceProvider;

        public DeploymentFileFactory(IServiceProvider serviceProvider)
		{
			this.serviceProvider = serviceProvider;
		}

        public DeploymentFile Create()
        {
            return this.serviceProvider.GetRequiredService<DeploymentFile>();
        }
    }
}

