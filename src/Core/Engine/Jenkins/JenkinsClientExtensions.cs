using System;
using JenkinsNET;
using System.Net.Http.Headers;
using System.Text;
using Microsoft.Extensions.Options;

namespace Modm.Engine.Jenkins
{
    public static class JenkinsClientExtensions
    {
        public static async Task<EngineStatus> GetJenkinsStatusAsync(this IJenkinsClient client)
        {
            EngineStatus status = new EngineStatus { EngineType = EngineType.Jenkins, IsHealthy = false, Version = "Unknown"};

            try
            {
                var jenkinsClient = client as JenkinsClient;
                if (jenkinsClient == null)
                {
                    return status;
                }

                var plainTextBytes = System.Text.Encoding.UTF8.GetBytes($"{jenkinsClient.UserName}:{jenkinsClient.Password}");
                var credential = Convert.ToBase64String(plainTextBytes);

                using (HttpClient httpClient = new HttpClient())
                {
                    // Add Authorization header with the API token
                    httpClient.DefaultRequestHeaders.Authorization =
                        new AuthenticationHeaderValue(
                            "Basic",
                            credential);

                    HttpResponseMessage response = await httpClient
                        .GetAsync(jenkinsClient.BaseUrl);

                    response.EnsureSuccessStatusCode();
                    status.IsHealthy = true;
                    string responseBody = await response.Content.ReadAsStringAsync();
                    string xJenkinsHeader = response.Headers.GetValues("X-Jenkins").FirstOrDefault();
                    status.Version = !string.IsNullOrEmpty(xJenkinsHeader) ? xJenkinsHeader : status.Version;
                    return status;
                }
            }
            catch (Exception ex)
            {
                return status;
            }
        }
    }
}

