// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
using Modm.Extensions;
using Modm.ServiceHost;

var host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureAppConfiguration((context, builder) =>
    {
        builder.AddEnvironmentVariables();
        builder.AddAppConfigurationSafely(context.HostingEnvironment);
    })
    .ConfigureServices((context, services) =>
    {
        services.AddServiceHost(context);
    })
    .Build();

await host.RunAsync();