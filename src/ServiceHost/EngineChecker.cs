// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;
using System.Timers;
using Microsoft.Extensions.Logging;
using Modm.Engine;

namespace Modm.ServiceHost
{
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
                var engineInfo = JsonSerializer.Deserialize<EngineInfo>(jsonResponse);
                if (engineInfo == null)
                {
                    this.logger.LogError($"Engine is not healthy. engineInfo is null.");
                    return false;
                }
                this.logger.LogInformation($"Engine status after deserialization: {engineInfo.IsHealthy}");
                return engineInfo.IsHealthy;
            }
            catch (Exception ex)
            {
                this.logger.LogError(ex, null);
                return false;
            }
        }
    }
}


