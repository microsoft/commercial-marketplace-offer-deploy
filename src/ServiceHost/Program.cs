using Modm.Azure;
using Modm.ServiceHost;

IHost host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureServices(services =>
    {
        services.AddHttpClient();
        services.AddSingleton<InstanceMetadataService>();
        services.AddHostedService<Startup>();
    })
    .Build();

await host.RunAsync();
