using System;
using System.Threading.Tasks;
using System.Text.Json;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
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

        [HttpGet("deployments")]
        public async Task<IActionResult> GetDeployments()
        {
            string backendUrl = this.configuration[BackendUrlSettingName];
            if (string.IsNullOrEmpty(backendUrl))
            {
                return BadRequest("Backend URL is not configured.");
            }
            try
            {
                var response = await client.GetAsync($"{backendUrl}/api/deployments");
                if (response.IsSuccessStatusCode)
                {
                    var content = await response.Content.ReadAsStringAsync();
                    var engineInfo = JsonSerializer.Deserialize<GetDeploymentResponse>(content, new JsonSerializerOptions { PropertyNameCaseInsensitive = true });
                    return Ok(engineInfo);
                }
                else
                {
                    return StatusCode((int)response.StatusCode, await response.Content.ReadAsStringAsync());
                }
            }
            catch (HttpRequestException ex)
            {
                // Log exception details (use ILogger for logging)
                return StatusCode(503, "Unable to reach the backend service."); // Service Unavailable
            }
            catch (JsonException ex)
            {
                // Log exception details (use ILogger for logging)
                return StatusCode(500, "Error parsing the response from the backend service."); // Internal Server Error
            }
            catch (Exception ex)
            {
                // Log exception details (use ILogger for logging)
                return StatusCode(500, "An unexpected error occurred."); // Internal Server Error
            }
        }

        [HttpGet("status")]
        public async Task<IActionResult> GetStatus()
        {
            string backendUrl = this.configuration[BackendUrlSettingName];
            if (string.IsNullOrEmpty(backendUrl))
            {
                return BadRequest("Backend URL is not configured.");
            }

            try
            {
                var response = await client.GetAsync($"{backendUrl}/api/status");
                if (response.IsSuccessStatusCode)
                {
                    var content = await response.Content.ReadAsStringAsync();
                    var engineInfo = JsonSerializer.Deserialize<EngineInfo>(content, new JsonSerializerOptions { PropertyNameCaseInsensitive = true });
                    return Ok(engineInfo);
                }
                else
                {
                    return StatusCode((int)response.StatusCode, await response.Content.ReadAsStringAsync());
                }
            }
            catch (HttpRequestException ex)
            {
                // Log exception details (use ILogger for logging)
                return StatusCode(503, "Unable to reach the backend service."); // Service Unavailable
            }
            catch (JsonException ex)
            {
                // Log exception details (use ILogger for logging)
                return StatusCode(500, "Error parsing the response from the backend service."); // Internal Server Error
            }
            catch (Exception ex)
            {
                // Log exception details (use ILogger for logging)
                return StatusCode(500, "An unexpected error occurred."); // Internal Server Error
            }
        }
    }
}