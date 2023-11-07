using Modm.ClientApp.Controllers;
using ClientApp.Security;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;
using System.Text;
using Microsoft.AspNetCore.Authorization;
using Microsoft.IdentityModel.Logging;

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

builder.Services.AddAuthentication(options =>
{
    options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
    options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
    options.DefaultScheme = JwtBearerDefaults.AuthenticationScheme;
}).AddJwtBearer(new JwtBearerConfigurator(builder.Configuration).Configure);

builder.Services.AddAuthorization();

builder.Services.AddControllersWithViews();
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

app.UseAuthentication();
app.UseAuthorization();


app.UseHttpsRedirection();
app.UseStaticFiles();

app.MapControllerRoute(
    name: "default",
    pattern: "{controller}/{action=Index}/{id?}"
).RequireAuthorization();

app.MapFallbackToFile("index.html");

app.Run();
