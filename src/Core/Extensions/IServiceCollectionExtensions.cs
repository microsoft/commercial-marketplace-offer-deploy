using System;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Modm.Artifacts;
using Modm.Engine;
using Modm.Engine.Jenkins.Client;

namespace Modm.Extensions
{
	public static class IServiceCollectionExtensions
	{

        /// <summary>
        /// Add the MODM deployment engine to the service collection
        /// </summary>
        /// <param name="services"></param>
        /// <returns></returns>
        /// <exception cref="NullReferenceException"></exception>
		public static IServiceCollection AddDeploymentEngine(this IServiceCollection services, IConfiguration configuration)
		{
            services.AddSingleton<ArtifactsDownloader>();

            services.AddSingleton<ApiTokenClient>();
            services.AddSingleton<JenkinsClientFactory>();

            services.AddSingleton<IJenkinsClient>(provider =>
            {
                var factory = provider.GetService<JenkinsClientFactory>();
                return factory == null ? throw new NullReferenceException("JenkinsClientFactory not configured") : factory.Create().GetAwaiter().GetResult();
            });

            services.AddSingleton<IDeploymentEngine, JenkinsDeploymentEngine>();


            //configuration
            services.Configure<JenkinsOptions>(configuration.GetSection(JenkinsOptions.ConfigSectionKey));

            services.AddSingletonHostedService<DeploymentMonitorService>();
            return services;
		}

        public static IServiceCollection AddSingletonHostedService<T>(this IServiceCollection services) where T : class, IHostedService
        {
            services.AddSingleton<T>().AddHostedService(p => p.GetRequiredService<T>());
            return services;
        }
    }
}

