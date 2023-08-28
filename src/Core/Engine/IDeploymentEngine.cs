namespace Modm.Engine
{
    public interface IDeploymentEngine
    {
        Task<StartDeploymentResult> StartAsync(string artifactsUri);
    }
}