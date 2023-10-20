using System;
using System.Net;
using Microsoft.Extensions.DependencyInjection;
using NSubstitute;

namespace Modm.Tests.Fakes
{
    public class FakeHttpMessageHandler : HttpMessageHandler
    {
        protected override Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellationToken)
        {
            return Task.FromResult(new HttpResponseMessage(HttpStatusCode.OK));
        }
    }

    public static class FakeHttpClient
    {
        public static IServiceCollection AddFakeHttpClient(this IServiceCollection services)
        {
            var httpClientFactory = Substitute.For<IHttpClientFactory>();
            var httpClient = new HttpClient(new FakeHttpMessageHandler()) { BaseAddress = new Uri("https://localhost") };
            httpClientFactory.CreateClient(Arg.Any<string>()).Returns(httpClient);


            services.AddSingleton(httpClient);
            services.AddSingleton(httpClientFactory);

            return services;
        }
    }
}

