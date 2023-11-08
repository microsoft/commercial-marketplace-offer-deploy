using Azure.Identity;
using Modm.Configuration;
using Modm.Extensions;
using Modm.ServiceHost;

var host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureAppConfiguration((context, builder) =>
    {
        builder.AddEnvironmentVariables();

        Console.WriteLine("Environment Name = " + context.HostingEnvironment.EnvironmentName);

        Console.WriteLine("IsPacker = " + context.HostingEnvironment.IsPacker());

        if (!context.HostingEnvironment.IsDevelopment() && !context.HostingEnvironment.IsPacker())
        {
            var provider = AppConfigurationEndpointProvider.New(context, builder.Build());

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
