using Modm.Azure;
using Modm.ServiceHost;

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
        services.AddSingleton<ArtifactsWatcher>();

        services.AddHostedService<ControllerService>();
        services.AddHostedService<ArtifactsWatcherService>();

        services.AddMediatR(c => c.RegisterServicesFromAssemblyContaining<ControllerService>());
    })
    .Build();

await host.RunAsync();
