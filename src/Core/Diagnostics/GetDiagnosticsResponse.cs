using System;
namespace Modm.Diagnostics
{
	public record GetDiagnosticsResponse
	{
		public string DeploymentEngine { get; set; }
	}
}

