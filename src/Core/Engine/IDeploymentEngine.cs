﻿using Modm.Deployments;

namespace Modm.Engine
{
    public interface IDeploymentEngine
    {
        Task<StartDeploymentResult> Start(string artifactsUri);

        /// <summary>
        /// Gets the status of the engine
        /// </summary>
        /// <returns></returns>
        Task<Deployment> Get();


        /// <summary>
        /// Get information about the engine
        /// </summary>
        /// <returns></returns>
        Task<EngineInfo> GetInfo();
    }
}