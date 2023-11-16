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

            if (deploymentType == DeploymentType.Arm)
			{
				return new ArmParametersFile(destinationDirectory);
			}	

			throw new ArgumentException($"Deployment type {deploymentType} not supported yet.");
		}
	}
}

