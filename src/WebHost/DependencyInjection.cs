using WebHost.Controllers;
using FluentValidation;
using Microsoft.Extensions.Azure;
using Azure.Identity;
using Modm.Extensions;
using Modm.Deployments;

namespace Modm.WebHost
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
		public static IServiceCollection AddWebHost(this IServiceCollection services, IConfiguration configuration, IHostEnvironment environment)
		{
            services.AddDefaultHttpClient();

            services.AddDeploymentEngine(configuration, environment);

            services.Configure<HostOptions>(hostOptions =>
            {
                hostOptions.BackgroundServiceExceptionBehavior = BackgroundServiceExceptionBehavior.Ignore;
            });

            services.AddControllersWithViews();
            services.AddAzureClients(clientBuilder =>
            {
                clientBuilder.AddArmClient(configuration.GetSection("Azure"));
                clientBuilder.UseCredential(new DefaultAzureCredential());
            });

            services.AddMediatR(c =>
            {
                c.RegisterServicesFromAssemblyContaining<DeploymentsController>();
                //c.AddOpenBehavior(typeof(LoggingBehaviour<,>));
                //c.AddOpenBehavior(typeof(ValidationBehavior<,>));
            });

            services.AddValidatorsFromAssemblyContaining<StartDeploymentRequestValidator>();


            return services;
		}
	}
}

