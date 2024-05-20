using System;

namespace Modm.Deployments
{
	public interface IDeploymentFileFactory
	{
        DeploymentFile Create();
    }
}

