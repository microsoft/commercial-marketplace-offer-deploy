using System;
namespace ClientApp.Cleanup
{
	public interface ICleanupService
	{
		Task<bool> CleanupInstallAsync(string resourceGroup);
	}
}

