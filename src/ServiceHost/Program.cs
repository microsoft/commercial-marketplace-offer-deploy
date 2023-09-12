using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Configuration;
using Modm;
using Modm.Azure;
using Modm.ServiceHost;

IConfiguration configuration = new ConfigurationBuilder()
    .SetBasePath(Directory.GetCurrentDirectory())
    .AddJsonFile("appsettings.json")
    .Build();

IHost host = Host.CreateDefaultBuilder(args)
    .UseSystemd()
    .ConfigureServices(services =>
    {
        services.AddHttpClient();
        services.AddSingleton<InstanceMetadataService>();
        services.AddSingleton<ManagedIdentityService>();
        services.AddHostedService<Startup>();

        services.AddSingleton<ArtifactsWatcher>(provider =>
        {
            var logger = provider.GetRequiredService<ILogger<ArtifactsWatcher>>();
            var httpClient = provider.GetRequiredService<HttpClient>();
            string artifactsFilePath = configuration.GetSection("ArtifactsWatcherSettings")["ArtifactsFilePath"];
            string statusEndpoint = configuration.GetSection("ArtifactsWatcherSettings")["DeploymentsEndpoint"];
            return new ArtifactsWatcher(httpClient, artifactsFilePath, statusEndpoint, logger);
        });
    })
    .Build();

await host.RunAsync();
