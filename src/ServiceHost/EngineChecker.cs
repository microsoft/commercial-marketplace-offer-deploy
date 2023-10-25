using System;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;
using System.Timers;

public class EngineChecker
{
    private readonly string BaseUrl;
    private readonly HttpClient httpClient;

    public EngineChecker(string baseUrl, HttpClient httpClient)
    {
        this.BaseUrl = baseUrl;
        this.httpClient = httpClient;
    }
  
    public async Task<bool> IsEngineHealthy()
    {
        try
        {
            var response = await httpClient.GetAsync($"{this.BaseUrl}/api/status");
            if (!response.IsSuccessStatusCode)
            {
                throw new Exception($"HTTP error! status: {response.StatusCode}");
            }

            var jsonResponse = await response.Content.ReadAsStringAsync();
            var statusData = JsonSerializer.Deserialize<StatusData>(jsonResponse);
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
