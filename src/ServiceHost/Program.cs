using Azure.Identity;
using Modm.Configuration;
using Modm.ServiceHost;

var host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureAppConfiguration((context, builder) =>
    {
        builder.AddEnvironmentVariables();

        if (!context.HostingEnvironment.IsDevelopment())
        {
            var provider = AppConfigurationEndpointProvider.FromHostBuilderContext(context);

            builder.AddAzureAppConfiguration(options =>
                  options.Connect(provider.Get(), new DefaultAzureCredential()));
        }
    })
    .ConfigureServices((context, services) =>
    {
        services.AddServiceHost(context);
    })
    .Build();

await host.RunAsync();
