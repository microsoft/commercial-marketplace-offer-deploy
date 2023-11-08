using Azure.Identity;
using Modm.Configuration;
using Modm.Extensions;
using Modm.ServiceHost;

var host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureAppConfiguration((context, builder) =>
    {
        builder.AddEnvironmentVariables();

        if (!context.HostingEnvironment.IsDevelopment())
        {
            try
            {
                var provider = AppConfigurationEndpointProvider.New(context, builder.Build());
                builder.AddAzureAppConfiguration(options => options.Connect(provider.Get(), new DefaultAzureCredential()));
            }
            catch
            {
            }
        }
    })
    .ConfigureServices((context, services) =>
    {
        services.AddServiceHost(context);
    })
    .Build();

await host.RunAsync();
