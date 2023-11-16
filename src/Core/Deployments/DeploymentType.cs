using System;
namespace Modm.Deployments
{
	/// <summary>
	/// Defines the types of deployments that MODM supports
	/// </summary>
	public static class DeploymentType
	{
		public static readonly string Arm = "arm";
        public static readonly string Terraform = "terraform";
    }
}

