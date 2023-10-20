using System.Net;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Extensions.Logging;

namespace Functions
{
    public class DashboardRedirect
    {
        public const string RedirectUrlAppSettingName = "RedirectUrl";
        public const string FallbackRedirectUrl = "https://portal.azure.com";

        private readonly ILogger logger;

        public DashboardRedirect(ILoggerFactory loggerFactory)
        {
            logger = loggerFactory.CreateLogger<DashboardRedirect>();
        }

        /// <summary>
        /// Performs a 302 redirect to create a known dashboard URL
        /// </summary>
        /// <param name="request"></param>
        /// <returns></returns>
        [Function("dashboard")]
        public async Task<HttpResponseData> Run([HttpTrigger(AuthorizationLevel.Anonymous, "get")] HttpRequestData request)
        {
            var redirectUrl = Environment.GetEnvironmentVariable(RedirectUrlAppSettingName, EnvironmentVariableTarget.Process) ?? string.Empty;

            if (string.IsNullOrEmpty(redirectUrl))
            {
                logger.LogError("The {settingName} returned an empty value. Serving up default /index.html", RedirectUrlAppSettingName);
                return await ReturnIndexPage(request);
            }

            logger.LogInformation("Redirecting request to: {url}", redirectUrl);

            var response = request.CreateResponse(HttpStatusCode.Found);
            response.Headers.Add("Location", redirectUrl);

            return response;
        }

        private static async Task<HttpResponseData> ReturnIndexPage(HttpRequestData request)
        {
            var page = Path.Combine(Environment.CurrentDirectory, "index.html");

            var response = request.CreateResponse(HttpStatusCode.OK);
            response.Headers.Add("Content-Type", "text/html; charset=utf-8");

            await response.WriteBytesAsync(await File.ReadAllBytesAsync(page));

            return response;
        }
    }
}
