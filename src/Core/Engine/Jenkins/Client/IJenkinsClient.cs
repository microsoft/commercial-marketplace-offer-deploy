using System;
using JenkinsNET;
using Modm.Engine.Jenkins.Model;

namespace Modm.Engine.Jenkins.Client
{
	/// <summary>
	/// Decorating client of the <see cref="JenkinsNET.IJenkinsClient"/> so we can include additional functionality
	/// to interact with Jenkins
	/// </summary>
	public interface IJenkinsClient : JenkinsNET.IJenkinsClient
    {
        /// <summary>
        /// Gets information about Hudson including the version
        /// </summary>
        /// <returns></returns>
        Task<JenkinsInfo> GetInfo();

        /// <summary>
        /// Gets the built-in node information
        /// </summary>
        /// <returns></returns>
        Task<MasterComputer> GetBuiltInNode();
	}
}

