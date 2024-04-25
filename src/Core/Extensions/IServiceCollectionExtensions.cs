using FluentValidation;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Modm.Packaging;
using Modm.Azure;
using Modm.Deployments;
using Modm.Engine;
using Modm.Jenkins.Client;
using Modm.Engine.Pipelines;
using Modm.Jenkins;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Modm.Security;

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
            services.AddScoped<IValidator<PackageFile>, PackageFileValidator>();
            services.AddSingleton<PackageFileFactory>();
            services.AddSingleton<ParametersFileFactory>();

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

            services.AddSingleton<IPackageDownloader, PackageDownloader>();

            services.AddSingleton<ApiTokenClient>();
            services.AddSingleton<JenkinsClientFactory>();
            services.AddSingleton<DeploymentFile>();
            services.AddSingleton<AuditFile>();
            services.AddSingleton<IDeploymentRepository, DefaultDeploymentRepository>();

            services.AddSingleton<IDeploymentEngine, JenkinsDeploymentEngine>();
            services.AddSingleton<DeploymentResourcesClient>();

            //configuration
            services.Configure<JenkinsOptions>(configuration.GetSection(JenkinsOptions.ConfigSectionKey));

            services.AddSingletonHostedService<JenkinsMonitorService>();
            services.AddSingletonHostedService<JenkinsReadinessService>();

            services.AddMediatR(c =>
            {
                c.RegisterServicesFromAssemblyContaining<IDeploymentEngine>();
            });

            services.AddPipeline<IPipeline<StartDeploymentRequest, StartDeploymentResult>, StartDeploymentRequestPipeline>(c => c.AddStartDeploymentRequestPipeline());
            services.AddPipeline<IPipeline<StartRedeploymentRequest, StartRedeploymentResult>, StartRedeploymentRequestPipeline>(c => c.AddStartRedeploymentRequestPipeline());

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
            services.AddHttpClient();
            return services;
        }

        /// <summary>
        /// Adds MODM jwt bearer token authentication
        /// </summary>
        /// <param name="services"></param>
        /// <param name="configuration"></param>
        /// <returns></returns>
        public static IServiceCollection AddJwtBearerAuthentication(this IServiceCollection services, IConfiguration configuration)
        {
            services.AddAuthentication(options =>
            {
                options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
                options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
                options.DefaultScheme = JwtBearerDefaults.AuthenticationScheme;
            }).AddJwtBearer(new JwtBearerConfigurator(configuration).Configure);

            services.AddAuthorization();

            return services;
        }
    }
}

