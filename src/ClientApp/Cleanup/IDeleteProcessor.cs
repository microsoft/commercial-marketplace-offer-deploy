using System;
namespace ClientApp.Cleanup
{
	public interface IDeleteProcessor
	{
		Task DeleteInstallResourcesAsync(string resourceGroup, CancellationToken cancellationToken);
	}
}

