using System;
namespace ClientApp.Cleanup
{
	public interface IDeleteProcessor
	{
		Task<bool> DeleteInstallResourcesAsync(string resourceGroup);
	}
}

