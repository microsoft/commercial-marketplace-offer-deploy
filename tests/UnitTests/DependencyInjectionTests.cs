using System;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using NSubstitute;
using Modm.ServiceHost;
using Xunit.Abstractions;
using System.ServiceProcess;
using Microsoft.Extensions.Configuration;
using Modm.WebHost;
using Modm.Jenkins.Client;
using Microsoft.Extensions.DependencyInjection.Extensions;

namespace Modm.Tests.UnitTests
{
    public class DependencyInjectionTests
    {
        private readonly ITestOutputHelper output;

        public DependencyInjectionTests(ITestOutputHelper output)
        {
            this.output = output;
        }

        [Fact]
        public void ServiceHost_should_wire_up_successfully()
        {
            var context = GetHostBuilderContext();
            var services = GetServiceCollection(context);

            services.AddServiceHost(context);

            VerifyServices(services);
        }

        [Fact]
        public void WebHost_should_wire_up_successfully()
        {
            var context = GetHostBuilderContext();
            var services = GetServiceCollection(context);

            services.AddWebHost(context.Configuration, context.HostingEnvironment);


            // mock jenkins client factory to disconnect from HTTP
            var factory = Substitute.For<JenkinsClientFactory>();
            factory.Create().ReturnsForAnyArgs(Substitute.For<IJenkinsClient>());
            services.Replace(new ServiceDescriptor(typeof(JenkinsClientFactory), factory));

            VerifyServices(services);
        }

        private void VerifyServices(IServiceCollection services)
        {
            var descriptors = services.Where(s => s.ImplementationType?.FullName?.Contains("Modm") ?? false).ToList();
            output.WriteLine($"service count = {descriptors.Count}");

            var provider = services.BuildServiceProvider();

            foreach (var service in descriptors)
            {
                output.WriteLine($"{service.ServiceType?.Name} => {service.ImplementationType?.Name}");

                Assert.NotNull(service.ServiceType);

                var exception = Record.Exception(() => provider.GetRequiredService(service.ServiceType));
                Assert.Null(exception);
            }
        }

        private static IServiceCollection GetServiceCollection(HostBuilderContext context)
        {
            var services = new ServiceCollection();
            services.AddSingleton(context.Configuration);

            return services;
        }

        private static HostBuilderContext GetHostBuilderContext()
        {
            var hostingEnvironment = Substitute.For<IHostEnvironment>();
            hostingEnvironment.EnvironmentName = "Development";

            var rootConfiguration = new ConfigurationBuilder()
                .AddInMemoryCollection(new List<KeyValuePair<string, string?>>()
                {
                    new KeyValuePair<string, string?>("Azure:DefaultSubscriptionId", Guid.NewGuid().ToString())
                })
                .Build();

            var context = new HostBuilderContext(new Dictionary<object, object>())
            {
                HostingEnvironment = hostingEnvironment,
                Configuration = rootConfiguration
            };

            return context;
        }
    }
}

