using System;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;
using System.Timers;

public class EngineChecker
{
    private readonly string baseUrl = "http://localhost:5000";
    private readonly HttpClient httpClient;
    ILogger<EngineChecker> logger;

    public EngineChecker(HttpClient httpClient, ILogger<EngineChecker> logger)
    {
        this.httpClient = httpClient;
        this.logger = logger;
    }
  
    public async Task<bool> IsEngineHealthy()
    {
        try
        {
            var response = await httpClient.GetAsync($"{this.baseUrl}/api/status");
            if (!response.IsSuccessStatusCode)
            {
                this.logger.LogError($"Engine is not healthy. Status code: {response.StatusCode}");
                return false;
            }

            var jsonResponse = await response.Content.ReadAsStringAsync();
            this.logger.LogInformation($"Engine status: {jsonResponse}");
            var statusData = JsonSerializer.Deserialize<StatusData>(jsonResponse);
            if (statusData == null)
            {
                this.logger.LogError($"Engine is not healthy. Status data is null.");
                return false;
            }
            this.logger.LogInformation($"Engine status after deserialization: {statusData.IsHealthy}");
            if (statusData == null)
            {
                return false;
            }

            return statusData.IsHealthy;
        }
        catch (Exception)
        {
            return false;
        }
    }

}

public class StatusData
{
    public bool IsHealthy { get; set; }
}
