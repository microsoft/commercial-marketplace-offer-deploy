namespace ClientApp.Cleanup
{
    public interface IDeleteProcessor
	{
		Task DeleteResourcesAsync(string resourceGroup, CancellationToken cancellationToken = default);
	}
}