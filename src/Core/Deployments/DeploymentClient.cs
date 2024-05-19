using System;
using System.Text.Json;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;

namespace Modm.Deployments
{
	public class DeploymentClient
	{
        private const string relativeUri = "api/deployments";
        private readonly HttpClient client;
        private readonly ILogger<DeploymentClient> logger;

        public DeploymentClient(HttpClient client, ILogger<DeploymentClient> logger)
		{
			this.client = client;
            this.logger = logger;
		}

		public async Task<GetDeploymentResponse> GetDeploymentInfo()
		{
            try
            {
                var response = await client.GetAsync(relativeUri);

                if (response.IsSuccessStatusCode)
                {
                    var content = await response.Content.ReadAsStringAsync();

                    JsonSerializerOptions serializerOptions = new()
                    {
                        PropertyNameCaseInsensitive = true,
                        PropertyNamingPolicy = JsonNamingPolicy.CamelCase
                    };

                    return JsonSerializer.Deserialize<GetDeploymentResponse>(content, serializerOptions);
                }
            }
            catch (HttpRequestException e)
            {
                const string message = "Unable to reach the backend service.";
                logger.LogError(e, message);
            }
            catch (JsonException e)
            {
                const string message = "Error parsing the response from the backend service.";
                logger.LogError(e, message);
            }
            catch (Exception e)
            {
                const string message = "An unexpected error occurred.";
                logger.LogError(e, message);
            }

            return null;
        }
    }
}

