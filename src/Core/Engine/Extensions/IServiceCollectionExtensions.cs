using System;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Modm.Artifacts;
using Modm.Engine.Jenkins.Client;

namespace Modm.Engine.Extensions
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

            services.AddScoped<IJenkinsClient>(provider =>
            {
                var factory = provider.GetService<JenkinsClientFactory>();
                return factory == null ? throw new NullReferenceException("JenkinsClientFactory not configured") : factory.Create().GetAwaiter().GetResult();
            });

            services.AddSingleton<IDeploymentEngine, JenkinsDeploymentEngine>();


            //configuration
            services.Configure<ArtifactsDownloadOptions>(configuration.GetSection(ArtifactsDownloadOptions.ConfigSectionKey));
            services.Configure<JenkinsOptions>(configuration.GetSection(JenkinsOptions.ConfigSectionKey));

            return services;
		}
	}
}

