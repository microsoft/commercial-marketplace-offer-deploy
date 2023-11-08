using Azure.Identity;
using System.Configuration;
using Modm.WebHost;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddWebHost(builder.Configuration, builder.Environment);

builder.Configuration.AddEnvironmentVariables();

builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowLocal", builder =>
    {
        builder.WithOrigins("https://localhost:44482");
    });
});

var appConfigEndpoint = builder.Configuration["Azure:AppConfigEndpoint"] ?? string.Empty;

Console.WriteLine(appConfigEndpoint);
Console.WriteLine(builder.Environment.IsDevelopment());

if (!string.IsNullOrEmpty(appConfigEndpoint))
{
    builder.Configuration.AddAzureAppConfiguration(options =>
      options.Connect(
          new Uri(appConfigEndpoint),
          new DefaultAzureCredential()));
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
