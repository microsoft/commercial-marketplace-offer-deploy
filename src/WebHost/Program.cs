using Azure.Identity;
using System.Configuration;
using Modm.WebHost;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddWebHost(builder.Configuration, builder.Environment);

builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowLocal", builder =>
    {
        builder.WithOrigins("https://localhost:44482");
    });
});

if (!builder.Environment.IsDevelopment())
{
    var appConfigEndpoint = builder.Configuration["Azure:AppConfigEndpoint"] ?? string.Empty;

    if (!string.IsNullOrEmpty(appConfigEndpoint))
    {
        builder.Configuration.AddAzureAppConfiguration(options =>
          options.Connect(
              new Uri(appConfigEndpoint),
              new ManagedIdentityCredential()));
    }
}

var app = builder.Build();

// Configure the HTTP request pipeline.
if (!app.Environment.IsDevelopment())
{
    // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
    app.UseHsts();
}

app.UseCors("AllowLocal");
app.UseHttpsRedirection();

app.MapControllerRoute(
    name: "default",
    pattern: "{controller}/{action=Index}/{id?}");

app.Run();
