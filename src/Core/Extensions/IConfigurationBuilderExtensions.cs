using System;
using Azure.Identity;
using Azure.ResourceManager;
using Microsoft.Extensions.Azure;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
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
        /// <param name="context"></param>
        /// <returns></returns>
        public static IConfigurationBuilder AddAppConfigurationIfExists(this IConfigurationBuilder builder, HostBuilderContext context)
        {
            var services = new ServiceCollection();
            var configuration = builder.Build();

            services.AddSingleton<IConfiguration>(configuration);
            services.AddDefaultHttpClient();

            services.AddAzureClients(clientBuilder =>
            {
                clientBuilder.AddArmClient(configuration.GetSection("Azure"));
                clientBuilder.UseCredential(new DefaultAzureCredential());
            });

            if (context.HostingEnvironment.IsDevelopment())
            {
                services.AddSingleton<IMetadataService, LocalMetadataService>();
            }
            else
            {
                services.AddSingleton<IMetadataService, DefaultMetadataService>();
            }

            services.AddSingleton(provider =>
            {
                return new AppConfigurationRegistration(
                    provider.GetRequiredService<IMetadataService>(),
                    provider.GetRequiredService<ArmClient>(),
                    builder);
            });

            var provider = services.BuildServiceProvider();
            var instance = provider.GetRequiredService<AppConfigurationRegistration>();

            instance.AddAppConfigurationIfExists();

            return builder;
        }
    }
}

