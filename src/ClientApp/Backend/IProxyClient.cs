using Microsoft.AspNetCore.Mvc;
using Microsoft.Azure.Management.ContainerInstance.Fluent;

namespace ClientApp.Backend
{
    public interface IProxyClient
    {
        Task<IActionResult> GetAsync<T>(string relativeUri);
        Task<IActionResult> PostAsync(string relativeUri, HttpContent content = default);
    }

    public static class ProxyClientType
    {
        public const string Http = "Http";
        public const string JsonFile = "JsonFile";


        public static bool IsHttp(string value)
        {
            return !string.IsNullOrEmpty(value) && value.Equals(Http, StringComparison.InvariantCultureIgnoreCase);
        }

        public static bool IsJsonFile(string value)
        {
            if (!string.IsNullOrEmpty(value))
            {
                return value.Contains(JsonFile, StringComparison.InvariantCultureIgnoreCase);
            }
            return false;
        }
    }
}

