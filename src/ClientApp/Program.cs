using Modm.ClientApp.Controllers;
using ClientApp.Security;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Logging;
using Modm.Security;
using Modm.Extensions;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton<AuthManager>();
builder.Services.AddSingleton<AdminCredentialsProvider>();

// configures the http client for the proxy controller to have requests proxied
builder.Services.AddHttpClient<ProxyController>().ConfigureHttpClient((provider, client) =>
{
    var backendUrl = provider.GetRequiredService<IConfiguration>()
                                .GetValue<string>(ProxyController.BackendUrlSettingName);
    client.BaseAddress = new Uri(backendUrl ?? string.Empty);
});

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
builder.Configuration.AddEnvironmentVariables();

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
