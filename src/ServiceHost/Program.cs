using Modm.Azure;
using Modm.Extensions;
using Modm.HttpClient;
using Modm.ServiceHost;
using Polly;
using Polly.Retry;
using Microsoft.Extensions.DependencyInjection;


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

        services.AddHttpClient(Constants.MODM)
        .AddTransientHttpErrorPolicy(builder => builder.WaitAndRetryAsync(new[]
        {
            TimeSpan.FromSeconds(1),
            TimeSpan.FromSeconds(5),
            TimeSpan.FromSeconds(10)
        }));

        services.AddSingletonHostedService<ControllerService>();
        services.AddSingletonHostedService<ArtifactsWatcherService>();
        services.AddSingletonHostedService<ManagedIdentityMonitorService>();

        services.AddMediatR(c => c.RegisterServicesFromAssemblyContaining<ControllerService>());
    })
    .Build();

await host.RunAsync();
