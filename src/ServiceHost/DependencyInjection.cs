using System;
using Modm.Azure;
using Modm.Deployments;
using Modm.Extensions;
using Polly;

namespace Modm.ServiceHost
{
    public static class ServiceCollectionExtensions
    {
        /// <summary>
        /// WebHost DI container registration
        /// </summary>
        /// <param name="services"></param>
        /// <param name="configuration"></param>
        /// <param name="environment"></param>
        /// <returns></returns>
		public static IServiceCollection AddServiceHost(this IServiceCollection services, HostBuilderContext context)
        {
            services.AddDefaultHttpClient();

            if (context.HostingEnvironment.IsDevelopment())
            {
                services.AddSingleton<IMetadataService, LocalMetadataService>();
                services.AddSingleton<IManagedIdentityService, LocalManagedIdentityService>();
            }
            else
            {
                services.AddSingleton<IMetadataService, DefaultMetadataService>();
                services.AddSingleton<IManagedIdentityService, DefaultManagedIdentityService>();
            }

            services.AddSingletonHostedService<ControllerService>();
            services.AddSingletonHostedService<PackageWatcherService>();
            services.AddSingletonHostedService<ManagedIdentityMonitorService>();

            services.AddMediatR(c => c.RegisterServicesFromAssemblyContaining<ControllerService>());

            return services;
        }
    }
}

