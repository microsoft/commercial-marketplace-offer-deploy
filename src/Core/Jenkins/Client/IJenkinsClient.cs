using Modm.Jenkins.Model;

namespace Modm.Jenkins.Client
{
    /// <summary>
    /// Decorating client of the <see cref="JenkinsNET.IJenkinsClient"/> so we can include additional functionality
    /// to interact with Jenkins
    /// </summary>
    public interface IJenkinsClient : IDisposable
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


        Task<string> GetBuildStatus(string jobName, int buildNumber);

        Task<bool> IsJobRunningOrWasAlreadyQueued(string jobName);

        Task<(int?, string)> Build(string jobName);

        Task<string> GetBuildLogs(string jobName, int buildNumber);

        Task<bool> IsBuilding(string jobName, int buildNumber, CancellationToken cancellationToken = default);
	}
}

