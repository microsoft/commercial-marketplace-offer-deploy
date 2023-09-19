using System;
namespace Modm.Deployments
{
	public class ParametersFileFactory
	{
		public ParametersFileFactory()
		{
		}

		public IDeploymentParametersFile Create(string deploymentType, string destinationDirectory)
		{
			if (deploymentType == DeploymentType.Terraform)
			{
				return new TerraformParametersFile(destinationDirectory);
			}

			throw new ArgumentException($"Deployment type {deploymentType} not supported yet.");
		}
	}
}

