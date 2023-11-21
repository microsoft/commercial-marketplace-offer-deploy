using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;
using JenkinsNET.Exceptions;
using JenkinsNET.Models;
using Microsoft.Extensions.Logging;
using Modm.Deployments;
using Modm.Jenkins.Model;
using Modm.Extensions;

namespace Modm.Jenkins.Client
{
    /// <summary>
    /// Our implementation against jenkins
    /// </summary>
    class JenkinsClient : IJenkinsClient
	{
        const string JenkinsVersionHeaderName = "X-Jenkins";

        private readonly HttpClient client;
        private readonly JenkinsNET.JenkinsClient jenkinsNetClient;
        private readonly JenkinsOptions options;
        private readonly ILogger<JenkinsClient> logger;

        public JenkinsClient(HttpClient client, JenkinsNET.JenkinsClient jenkinsNetClient, ILogger<JenkinsClient> logger, JenkinsOptions options)
        {
            this.client = client;
            this.jenkinsNetClient = jenkinsNetClient;
            this.logger = logger;
            this.options = options;
        }

        public async Task<JenkinsInfo> GetInfo()
        {
            
            var version = "";
            var hudson = await Send<Hudson>(HttpMethod.Get, "api/json", response =>
            {
                version = response.Headers.GetValues(JenkinsVersionHeaderName).FirstOrDefault();
                response.EnsureSuccessStatusCode();
            });

            return new JenkinsInfo { Hudson = hudson, Version = version };
        }

        public async Task<MasterComputer> GetBuiltInNode()
        {
            return await Send<MasterComputer>(HttpMethod.Get, "computer/(built-in)/api/json", response => response.EnsureSuccessStatusCode());
        }

        /// <summary>
        /// Helper method that sends an HTTP request to Jenkins using the method and relative Uri with an optional http response message handler,
        /// then deserializes the JSON response
        /// </summary>
        /// <typeparam name="T">The type of response to deserialize</typeparam>
        /// <param name="method"></param>
        /// <param name="relativeUri"></param>
        /// <param name="handler"></param>
        /// <returns></returns>
        private async Task<T> Send<T>(HttpMethod method, string relativeUri, Action<HttpResponseMessage> handler = null)
        {
            using var request = GetHttpRequest(method, relativeUri);
            var response = await client.SendAsync(request);

            if (handler != null)
            {
                handler(response);
            }

            return await Deserialize<T>(response);
        }

        private HttpRequestMessage GetHttpRequest(HttpMethod method, string relativeUri)
        {
            var requestUri = new Uri(options.BaseUrl).Append(relativeUri).AbsoluteUri;
            var request = new HttpRequestMessage(method, requestUri);

            if (String.IsNullOrEmpty(options.ApiToken))
            {
                options.ApiToken = this.jenkinsNetClient.ApiToken;
            }

            request.Headers.Authorization = GetAuthenticationHeader(options);

            return request;
        }

        /// <summary>
        /// gets an authentication header using the username and api token
        /// </summary>
        /// <param name="options"></param>
        /// <returns></returns>
        private static AuthenticationHeaderValue GetAuthenticationHeader(JenkinsOptions options)
        {
            var plainTextBytes = System.Text.Encoding.UTF8.GetBytes($"{options.UserName}:{options.ApiToken}");
            var credential = Convert.ToBase64String(plainTextBytes);

            return new AuthenticationHeaderValue("Basic", credential);
        }

        /// <summary>
        /// deserializes the response as JSON
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="response"></param>
        /// <returns></returns>
        private static async Task<T> Deserialize<T>(HttpResponseMessage response)
        {
            return await JsonSerializer.DeserializeAsync<T>(response.Content.ReadAsStream());
        }

        public async Task<string> GetBuildStatus(string jobName, int buildNumber)
        {
            var status = DeploymentStatus.Undefined;

            try
            {
                this.logger.LogInformation($"Inside GetBuildStatus prior to calling jenkinsNetClient.Builds.GetAsync with jobName:{jobName} and buildNumber:{buildNumber.ToString()}");
                var build = await jenkinsNetClient.Builds.GetAsync<JenkinsBuildBase>(jobName, buildNumber.ToString());
                if (build != null && !string.IsNullOrEmpty(build.Result))
                {
                    status = build.Result.ToLower();
                }
                
            }
            catch(Exception ex) // unfortunately we have to just catch all since the jenkinsNet client will throw if it doesn't exist
            {
                this.logger.LogError(ex.Message);
            }

            return status;
        }

        public async Task<int?> Build(string jobName)
        {
            var (id, status) = await ExecuteBuild(jobName);
            if (!id.HasValue)
            {
                bool isBuildStarted = false;
                int maxAttempts = 20;
                int attempt = 0;
                int delay = 5000;

                while (!isBuildStarted && attempt < maxAttempts)
                {
                    attempt++;
                    this.logger.LogInformation($"Waiting for build to start. Attempt: {attempt}");

                    await Task.Delay(delay);

                    id = await GetLastBuildNumberAsync(jobName);

                    if (id.HasValue)
                    {
                        this.logger.LogInformation("id.HasValue was true");
                        status = await GetBuildStatus(jobName, id.Value);
                        this.logger.LogInformation($"returned the following status - {status}");
                    }

                    isBuildStarted = (id.HasValue);
                }

                if (!isBuildStarted)
                {
                    this.logger.LogInformation("Build did not start after polling attempts.");
                }
            }

            return id;
        }

        private async Task<(int?, string)> ExecuteBuild(string jobName)
        {
            var response = await jenkinsNetClient.Jobs.BuildAsync(jobName);
            var queueId = response.GetQueueItemNumber().GetValueOrDefault(0);

            var failedToEnqueue = queueId == 0;

            if (failedToEnqueue)
            {
                return (null, DeploymentStatus.Undefined);
            }

            var queueItem = await jenkinsNetClient.Queue.GetItemAsync(queueId);
            var buildNumber = queueItem?.Executable?.Number;

            if (!buildNumber.HasValue)
            {
                return (null, DeploymentStatus.Undefined);
            }

            return (buildNumber.Value, DeploymentStatus.Running);
        }

        public async Task<int?> GetLastBuildNumberAsync(string jobName)
        {
            try
            {
                
                var lastBuild = await jenkinsNetClient.Builds.GetAsync<JenkinsBuildBase>(jobName, "lastBuild");

                if (lastBuild != null)
                {
                    return lastBuild.Number;
                }
            }
            catch (Exception ex)
            {
                logger.LogError(ex, $"Failed to fetch the last build number for job {jobName}.");
            }

            return null;
        }


        public async Task<bool> IsJobRunningOrWasAlreadyQueued(string jobName)
        {
            try
            {
                var queueItems = await jenkinsNetClient.Queue.GetAllItemsAsync();

                if (queueItems.Any(item => item.Task?.Name == jobName))
                {
                    logger.LogDebug($"Deployment {jobName} is already in the queue.");
                    return false;
                }

                var lastBuild = await jenkinsNetClient.Builds.GetAsync<JenkinsBuildBase>(jobName, "lastBuild");

                if (lastBuild == null)
                {
                    logger.LogWarning($"Unable to fetch last build details for {jobName}");
                    return true;
                }

                if (lastBuild.Building.GetValueOrDefault(false) || !string.IsNullOrEmpty(lastBuild.Result))
                {
                    logger.LogDebug($"Deployment {jobName} either is currently building or already completed.");
                    return false;
                }

                return true;
            }
            catch (JenkinsJobGetBuildException)
            {
                logger.LogWarning($"No previous builds found for {jobName} due to a 404 response. Assuming the job is startable.");
                return true;
            }
            catch (Exception ex)
            {
                logger.LogError(ex, $"Failed to check if the deployment is startable for {jobName}");
                return false;
            }
        }

        public async Task<string> GetBuildLogs(string jobName, int buildNumber)
        {
            try
            {
                var text = await jenkinsNetClient.Builds.GetConsoleTextAsync(jobName, buildNumber.ToString());

                if (string.IsNullOrEmpty(text))
                {
                    var bytes = Encoding.Default.GetBytes(text);
                    text = Encoding.UTF8.GetString(bytes);
                }
                return text;
            }
            catch
            {
            }
            return string.Empty;
        }

        public async Task<bool> IsBuilding(string jobName, int buildNumber, CancellationToken cancellationToken = default)
        {
            try
            {
                var build = await jenkinsNetClient.Builds.GetAsync<JenkinsBuildBase>(jobName, buildNumber.ToString(), cancellationToken);
                if (build == null)
                {
                    return false;
                }

                return build.Building.GetValueOrDefault(false);
            }
            catch (Exception ex)
            {
                logger.LogError(ex, "Exception thrown while trying to get build status");
            }

            return false;
        }

        public void Dispose()
        {
        }
    }
}

