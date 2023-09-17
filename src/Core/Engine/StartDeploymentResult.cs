using System;
namespace Modm.Engine
{
	public record StartDeploymentResult
	{
		public int Id { get; set; }
		public string Status { get; set; }
	}
}

