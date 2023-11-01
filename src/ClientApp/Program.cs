
using System.Net;
using Modm.ClientApp.Controllers;

var builder = WebApplication.CreateBuilder(args);

// configures the http client for the proxy controller to have requests proxied
builder.Services.AddHttpClient<ProxyController>().ConfigureHttpClient((provider, client) =>
{
    var backendUrl = provider.GetRequiredService<IConfiguration>()
                                .GetValue<string>(ProxyController.BackendUrlSettingName);
    client.BaseAddress = new Uri(backendUrl ?? string.Empty);
});

builder.Services.AddControllersWithViews();

builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowLocal", builder =>
    {
        builder.WithOrigins("https://localhost:44482");
    });
});

builder.Configuration.AddEnvironmentVariables();

var app = builder.Build();

if (!app.Environment.IsDevelopment())
{
    // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
    app.UseHsts();
}

app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseCors("AllowLocal");

app.MapControllerRoute(
    name: "default",
    pattern: "{controller}/{action=Index}/{id?}");

app.MapFallbackToFile("index.html");

app.Run();
