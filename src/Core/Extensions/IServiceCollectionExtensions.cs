using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Modm.Artifacts;
using Modm.Azure;
using Modm.Deployments;
using Modm.Engine;
using Modm.Engine.Jenkins.Client;
using Modm.Engine.Pipelines;
using Modm.Http;
using Polly;

namespace Modm.Extensions
{
    public static class IServiceCollectionExtensions
	{
        /// <summary>
        /// Adds a pipeline that's implemented using <see cref="MediatR"/>
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="services"></param>
        /// <returns></returns>
        public static IServiceCollection AddPipeline<TService, TImplementation>(this IServiceCollection services, Action<MediatRServiceConfiguration> configure)
            where TService : class
            where TImplementation : class, TService
        {
            services.AddSingleton<TService, TImplementation>();
            services.AddMediatR(configure);
            return services;
        }

        /// <summary>
        /// Add the MODM deployment engine to the service collection
        /// </summary>
        /// <param name="services"></param>
        /// <returns></returns>
        /// <exception cref="NullReferenceException"></exception>
		public static IServiceCollection AddDeploymentEngine(this IServiceCollection services, IConfiguration configuration, IHostEnvironment environment)
		{
            if (environment.IsDevelopment())
            {
                services.AddSingleton<IMetadataService, LocalMetadataService>();
                services.AddSingleton<IManagedIdentityService, LocalManagedIdentityService>();
            }
            else
            {
                services.AddSingleton<IMetadataService, DefaultMetadataService>();
                services.AddSingleton<IManagedIdentityService, DefaultManagedIdentityService>();
            }

            services.AddSingleton<ArtifactsDownloader>();

            services.AddSingleton<ApiTokenClient>();
            services.AddSingleton<JenkinsClientFactory>();
            services.AddScoped<DeploymentFile>();

            services.AddSingleton<IJenkinsClient>(provider =>
            {
                var factory = provider.GetService<JenkinsClientFactory>();
                return factory == null ? throw new NullReferenceException("JenkinsClientFactory not configured") : factory.Create().GetAwaiter().GetResult();
            });

            services.AddSingleton<IDeploymentEngine, JenkinsDeploymentEngine>();
            services.AddSingleton<DeploymentResourcesClient>();

            //configuration
            services.Configure<JenkinsOptions>(configuration.GetSection(JenkinsOptions.ConfigSectionKey));

            services.AddSingletonHostedService<JenkinsMonitorService>();

            services.AddMediatR(c =>
            {
                c.RegisterServicesFromAssemblyContaining<IDeploymentEngine>();
            });

            services.AddPipeline<IPipeline<StartDeploymentRequest, StartDeploymentResult>, StartDeploymentRequestPipeline>(c => c.AddStartDeploymentRequestPipeline());

            return services;
		}

        public static IServiceCollection AddSingletonHostedService<T>(this IServiceCollection services) where T : class, IHostedService
        {
            services.AddSingleton<T>().AddHostedService(p => p.GetRequiredService<T>());
            return services;
        }

        /// <summary>
        /// Adds our default, configured httpclient
        /// </summary>
        /// <param name="services"></param>
        /// <returns></returns>
        public static IServiceCollection AddDefaultHttpClient(this IServiceCollection services)
        {
            services.AddHttpClient(HttpConstants.DefaultHttpClientName)
                .AddTransientHttpErrorPolicy(builder => builder.WaitAndRetryAsync(new[]
                {
                            TimeSpan.FromSeconds(1),
                            TimeSpan.FromSeconds(5),
                            TimeSpan.FromSeconds(10)
                }));
            return services;
        }
    }
}

