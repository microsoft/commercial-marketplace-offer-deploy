using System;
using Azure.Identity;
using Azure.ResourceManager;
using Microsoft.Extensions.Azure;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Modm.Azure;
using Modm.Configuration;

namespace Modm.Extensions
{
	public static class IConfigurationBuilderExtensions
	{
        /// <summary>
        /// Will safely register the app configuration with the configuration against the builder IF
        /// the Azure resource actually exists
        /// </summary>
        /// <param name="environment"></param>
        /// <returns></returns>
        public static IConfigurationBuilder AddAppConfigurationSafely(this IConfigurationBuilder builder, IHostEnvironment environment)
        {
            var registrar = BuildRegistrar(builder, environment);
            registrar.AddAppConfigurationIfExists();

            return builder;
        }

        public static AppConfigurationRegistrar BuildRegistrar(IConfigurationBuilder builder, IHostEnvironment environment)
        {
            var services = new ServiceCollection();
            var configuration = builder.Build();

            services.AddLogging(loggingBuilder =>
            {
                loggingBuilder.AddConsole();
            });

            services.AddSingleton<IConfiguration>(configuration);
            services.AddDefaultHttpClient();

            services.AddAzureClients(clientBuilder =>
            {
                clientBuilder.AddArmClient(configuration.GetSection("Azure"));
                clientBuilder.UseCredential(new DefaultAzureCredential());
            });

            if (environment.IsDevelopment())
            {
                services.AddSingleton<IMetadataService, LocalMetadataService>();
            }
            else
            {
                services.AddSingleton<IMetadataService, DefaultMetadataService>();
            }

            services.AddSingleton(provider =>
            {
                return new AppConfigurationRegistrar(
                    provider.GetRequiredService<IMetadataService>(),
                    provider.GetRequiredService<ArmClient>(),
                    builder,
                    provider.GetRequiredService<ILogger<AppConfigurationRegistrar>>());
            });

            var provider = services.BuildServiceProvider();
            var instance = provider.GetRequiredService<AppConfigurationRegistrar>();

            return instance;
        }
    }
}