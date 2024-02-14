// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using Microsoft.AspNetCore.Mvc;
using System.Text.Json;

namespace ClientApp.Backend
{
    public class HttpProxyClient : IProxyClient
    {
        public const string BackendUrlSettingName = "BackendUrl";
        private readonly HttpClient client;
        private readonly ILogger<HttpProxyClient> logger;
        private readonly static JsonSerializerOptions serializerOptions = new()
        { 
            PropertyNameCaseInsensitive = true, 
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase 
        };

        public HttpProxyClient(HttpClient client, ILogger<HttpProxyClient> logger)
        {
            this.client = client;
            this.logger = logger;
        }

        public async Task<IActionResult> PostAsync(string relativeUri, HttpContent content = default)
        {
            try
            {
                var response = await client.PostAsync(relativeUri, content);
                return new ContentResult
                {
                    Content = await response.Content.ReadAsStringAsync(),
                    ContentType = "application/json",
                    StatusCode = (int)response.StatusCode
                };
            }
            catch (HttpRequestException e)
            {
                const string message = "Unable to reach the backend service.";
                logger.LogError(e, message);
                return StatusCode(503, message);
            }
            catch (Exception e)
            {
                const string message = "An unexpected error occurred.";
                logger.LogError(e, message);
                return StatusCode(500, message);
            }
        }

        public async Task<IActionResult> GetAsync<T>(string relativeUri)
        {
            try
            {
                var response = await client.GetAsync(relativeUri);

                if (response.IsSuccessStatusCode)
                {
                    var value = await DeserializeResponse<T>(response);
                    logger.LogTrace("Response from backend: {value}", JsonSerializer.Serialize(value));

                    return new OkObjectResult(value);
                }

                return StatusCode((int)response.StatusCode, await response.Content.ReadAsStringAsync());
            }
            catch (HttpRequestException e)
            {
                const string message = "Unable to reach the backend service.";
                logger.LogError(e, message);
                return StatusCode(503, message);
            }
            catch (JsonException e)
            {
                const string message = "Error parsing the response from the backend service.";
                logger.LogError(e, message);
                return StatusCode(500, message);
            }
            catch (Exception e)
            {
                const string message = "An unexpected error occurred.";
                logger.LogError(e, message);
                return StatusCode(500, message);
            }
        }

        private static async Task<T> DeserializeResponse<T>(HttpResponseMessage response)
        {
            var content = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<T>(content, serializerOptions);
        }

        private static ObjectResult StatusCode(int statusCode, string message)
        {
            return new ObjectResult(message)
            {
                StatusCode = statusCode
            };
        }
    }
}

