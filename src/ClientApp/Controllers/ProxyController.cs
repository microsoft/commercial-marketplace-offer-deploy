using System;
using System.Threading.Tasks;
using System.Text.Json;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using System.Net.Http.Headers;
using Modm.Deployments;
using Modm.Diagnostics;
using Modm.Engine;

namespace Modm.ClientApp.Controllers
{
    [Route("api")]
    [ApiController]
    [Authorize]
    public class ProxyController : ControllerBase
    {
        public const string BackendUrlSettingName = "BackendUrl";
        private readonly HttpClient client;
        private readonly IConfiguration configuration;

        public ProxyController(HttpClient client, IConfiguration configuration)
        {
            this.client = client;
            this.configuration = configuration;
        }

        private string RetrieveJwtToken()
        {
            if (HttpContext.Request.Headers.TryGetValue("Authorization", out var authHeader))
            {
                var token = authHeader.ToString().Split(' ').Last();
                return token;
            }

            return null;
        }

        private async Task<IActionResult> GetFromBackendService<T>(string endpoint)
        {
            string backendUrl = this.configuration[BackendUrlSettingName];
            if (string.IsNullOrEmpty(backendUrl))
            {
                return BadRequest("Backend URL is not configured.");
            }

            var token = RetrieveJwtToken();
            if (token == null)
            {
                return Unauthorized("JWT Token is missing");
            }

            client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);

            try
            {
                var response = await client.GetAsync($"{backendUrl}/api/{endpoint}");
                if (response.IsSuccessStatusCode)
                {
                    var content = await response.Content.ReadAsStringAsync();
                    var result = JsonSerializer.Deserialize<T>(content, new JsonSerializerOptions { PropertyNameCaseInsensitive = true });
                    return Ok(result);
                }
                else
                {
                    return StatusCode((int)response.StatusCode, await response.Content.ReadAsStringAsync());
                }
            }
            catch (HttpRequestException)
            {
                return StatusCode(503, "Unable to reach the backend service.");
            }
            catch (JsonException)
            {
                return StatusCode(500, "Error parsing the response from the backend service.");
            }
            catch (Exception)
            {
                return StatusCode(500, "An unexpected error occurred.");
            }
        }

        [HttpPost]
        [Route("resources/{resourceGroupName}/deletemodmresources")]
        public async Task<IActionResult> DeleteResourcesWithTagAsync([FromRoute] string resourceGroupName)
        {
            string backendUrl = this.configuration[BackendUrlSettingName];
            if (string.IsNullOrEmpty(backendUrl))
            {
                return BadRequest("Backend URL is not configured.");
            }

            var token = RetrieveJwtToken();
            if (token == null)
            {
                return Unauthorized("JWT Token is missing");
            }

            client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);

            try
            {
                var response = await client.PostAsync($"{backendUrl}/api/resources/{resourceGroupName}/deletemodmresources", null);

                var content = await response.Content.ReadAsStringAsync();
                return new ContentResult
                {
                    Content = content,
                    ContentType = "application/json",
                    StatusCode = (int)response.StatusCode
                };
            }
            catch (HttpRequestException)
            {
                return StatusCode(503, "Unable to reach the backend service.");
            }
            catch (Exception)
            {
                return StatusCode(500, "An unexpected error occurred.");
            }
        }


        [HttpGet("deployments")]
        public Task<IActionResult> GetDeployments()
        {
            return GetFromBackendService<GetDeploymentResponse>("deployments");
        }

        [HttpGet("diagnostics")]
        public Task<IActionResult> GetDiagnostics()
        {
            return GetFromBackendService<GetDiagnosticsResponse>("diagnostics");
        }

        [HttpGet("status")]
        public Task<IActionResult> GetStatus()
        {
            return GetFromBackendService<EngineInfo>("status");
        }
    }
}
