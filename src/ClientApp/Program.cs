using Modm.ClientApp.Controllers;
using ClientApp.Security;
using Microsoft.IdentityModel.Logging;
using Modm.Security;
using Modm.Extensions;
using ClientApp.Backend;
using Modm.Azure;
using MediatR;
using Microsoft.Azure.Management.Storage.Fluent.Models;
using System.Web.Services.Description;
using Microsoft.Extensions.Azure;
using Azure.Identity;
using Modm.Azure.Notifications;
using ClientApp.Notifications;
using ClientApp;
using System.Configuration;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton<AuthManager>();
builder.Services.AddSingleton<AdminCredentialsProvider>();

// configures the http client for the proxy controller to have requests proxied
builder.Services.AddHttpClient<ProxyController>().ConfigureHttpClient((provider, client) =>
{
    var backendUrl = provider.GetRequiredService<IConfiguration>()
                                .GetValue<string>(ProxyClientFactory.BackendUrlSettingName);
    client.BaseAddress = new Uri(backendUrl ?? string.Empty);
});

builder.Services.AddHttpClient<Modm.Deployments.DeploymentClient>().ConfigureHttpClient((provider, client) =>
{
    var backendUrl = provider.GetRequiredService<IConfiguration>()
                              .GetValue<string>(ProxyClientFactory.BackendUrlSettingName); 
    client.BaseAddress = new Uri(backendUrl ?? string.Empty);
});

builder.Services.AddSingleton<ProxyClientFactory>();

builder.Services.AddJwtBearerAuthentication(builder.Configuration);

builder.Services.AddControllersWithViews();
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowLocal", builder =>
    {
        builder
            .WithOrigins("https://localhost:44482", "https://localhost:7153", "http://localhost:5207", "http://localhost:44482")
            .AllowAnyMethod()
        .AllowAnyHeader();
        
    });
});

var dataDirectory = builder.Environment.IsDevelopment()
        ? builder.Configuration["DataDirectory"]
        : builder.Configuration["AppServiceDataDirectory"];

builder.Services.Configure<DeleteServiceOptions>(options => options.DataDirectory = dataDirectory);

builder.Services.Configure<HostOptions>(hostOptions =>
{
    hostOptions.BackgroundServiceExceptionBehavior = BackgroundServiceExceptionBehavior.Ignore;
});

builder.Services.AddAzureClients(clientBuilder =>
{
    clientBuilder.AddArmClient(builder.Configuration.GetSection("Azure"));
    clientBuilder.UseCredential(new DefaultAzureCredential());
});

builder.Services.AddSingleton<IAzureResourceManager, AzureResourceManager>();

builder.Services.AddSingletonHostedService<DeleteService>();
builder.Services.AddSingletonHostedService<AzureDeploymentCleanupService>();


builder.Services.AddMediatR(c =>
{
    c.RegisterServicesFromAssemblyContaining<DeleteInitiated>();
});



builder.Configuration.AddEnvironmentVariables();
builder.Configuration.AddAppConfigurationSafely(builder.Environment);

var app = builder.Build();

if (!app.Environment.IsDevelopment())
{
    // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
    app.UseHsts();
}
else
{
    IdentityModelEventSource.ShowPII = true;
}

app.UseHttpsRedirection();
app.UseStaticFiles();

app.UseRouting();

app.UseCors("AllowLocal");

app.UseAuthentication();
app.UseAuthorization();

app.MapControllerRoute(
    name: "default",
    pattern: "{controller}/{action=Index}/{id?}"
).RequireAuthorization();

app.MapFallbackToFile("index.html");

app.Run();
