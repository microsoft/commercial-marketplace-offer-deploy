using System.Net;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Extensions.Logging;

namespace Functions
{
    public class DashboardRedirect
    {
        private readonly ILogger _logger;

        public DashboardRedirect(ILoggerFactory loggerFactory)
        {
            _logger = loggerFactory.CreateLogger<DashboardRedirect>();
        }

        [Function("DashboardRedirect")]
        public HttpResponseData Run([HttpTrigger(AuthorizationLevel.Anonymous, "get", "post")] HttpRequestData req)
        {
            _logger.LogInformation("C# HTTP trigger function processed a request.");

            var response = req.CreateResponse(HttpStatusCode.Found);
            response.Headers.Add("Location", "http://www.microsoft.com");

            return response;
        }
    }
}
