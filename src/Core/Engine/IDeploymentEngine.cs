using Modm.Deployments;

namespace Modm.Engine
{
    public interface IDeploymentEngine
    {
        Task<StartDeploymentResult> Start(StartDeploymentRequest request, CancellationToken cancellationToken);

        /// <summary>
        /// Gets the status of the engine
        /// </summary>
        /// <returns></returns>
        Task<Deployment> Get();

        /// <summary>
        /// Gets the status of the engine
        /// </summary>
        /// <returns></returns>
        Task<string> GetLogs();


        /// <summary>
        /// Get information about the engine
        /// </summary>
        /// <returns></returns>
        Task<EngineInfo> GetInfo();
    }
}