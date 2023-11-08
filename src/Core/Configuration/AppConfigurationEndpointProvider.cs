using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Modm.Azure;
using Modm.Extensions;

namespace Modm.Configuration
{
    /// <summary>
    /// used for bootstrapping the appconfiguration by providing the appconfigstore endpoint
    /// using the Imds
    /// </summary>
    public class AppConfigurationEndpointProvider
	{
        private readonly IMetadataService metadataService;

        public AppConfigurationEndpointProvider(IMetadataService metadataService)
		{
            this.metadataService = metadataService;
        }

        /// <summary>
        /// Creates a new instance of the provider using an instance of <see cref="HostBuilderContext"/>
        /// </summary>
        /// <param name="context"></param>
        /// <returns></returns>
		public static AppConfigurationEndpointProvider FromHostBuilderContext(HostBuilderContext context)
		{
			var services = new ServiceCollection();

            services.AddSingleton(context.Configuration);
            services.AddDefaultHttpClient();

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
                return new AppConfigurationEndpointProvider(provider.GetRequiredService<IMetadataService>());
            });

            var provider = services.BuildServiceProvider();



            return provider.GetRequiredService<AppConfigurationEndpointProvider>();
        }

		public Uri Get()
		{
            var metadata = metadataService.GetAsync().GetAwaiter().GetResult();
            var resource = new AppConfigurationResource(metadata.ResourceGroupId);

            return resource.Uri;
        }
	}
}

