using Modm.ServiceHost;

var host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureAppConfiguration(builder =>
    {
        builder.AddEnvironmentVariables();
    })
    .ConfigureServices((context, services) => services.AddServiceHost(context))
    .Build();

await host.RunAsync();
