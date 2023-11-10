using System;
namespace ClientApp.Backend
{
	/// <summary>
	/// The API routes for the backend of MODM
	/// </summary>
	public static class Routes
	{
		public const string DeleteInstallerFormat = "api/resources/{0}/deletemodmresources";
        public const string GetDeployments = "api/deployments";
		public const string GetDiagnostics = "api/diagnostics";
		public const string GetStatus = "api/status";
    }
}

