using Modm.Extensions;
using Modm.ServiceHost;

var host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureAppConfiguration((context, builder) =>
    {
        builder.AddEnvironmentVariables();
        builder.AddAppConfigurationIfExists(context);
    })
    .ConfigureServices((context, services) =>
    {
        services.AddServiceHost(context);
    })
    .Build();

await host.RunAsync();