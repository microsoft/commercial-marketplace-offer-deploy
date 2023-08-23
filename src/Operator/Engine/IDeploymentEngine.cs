namespace Operator.Engine
{
    public interface IDeploymentEngine
    {
        Task<int> StartAsync(string artifactsUri);
    }
}