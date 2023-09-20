using Modm.Azure;
using Modm.Extensions;
using Modm.ServiceHost;
using Polly;
using Polly.Retry;



IHost host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureAppConfiguration(builder =>
    {
        builder.AddEnvironmentVariables();
    })
    .ConfigureServices((context, services) =>
    {
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

        services.AddHttpClient();

        services.AddSingletonHostedService<ControllerService>();
        services.AddSingletonHostedService<ArtifactsWatcherService>();
        services.AddSingletonHostedService<ManagedIdentityMonitorService>();

        services.AddMediatR(c => c.RegisterServicesFromAssemblyContaining<ControllerService>());
    })
    .Build();

await host.RunAsync();
