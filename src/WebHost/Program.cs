using Modm;
using Modm.Artifacts;
using Modm.Engine;
using WebHost.Deployments;
using FluentValidation;
using Microsoft.Extensions.Azure;
using Azure.Identity;
using Modm.Engine.Extensions;
using Modm.Deployments;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddHttpClient();

builder.Services.AddDeploymentEngine(builder.Configuration);
builder.Services.AddSingleton<ArtifactsDownloader>();

builder.Services.AddScoped<IValidator<CreateDeploymentRequest>, CreateDeploymentRequestValidator>();

// Add services to the container.
builder.Services.AddHostedService<Worker>();
builder.Services.AddControllersWithViews();

// azure configuration
builder.Services.AddAzureClients(clientBuilder =>
{
    clientBuilder.AddArmClient(builder.Configuration.GetSection("Azure"));
    clientBuilder.UseCredential(new DefaultAzureCredential());
});



var app = builder.Build();

// Configure the HTTP request pipeline.
if (!app.Environment.IsDevelopment())
{
    // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
    app.UseHsts();
}

app.UseHttpsRedirection();
app.UseStaticFiles();

app.MapControllerRoute(
    name: "default",
    pattern: "{controller}/{action=Index}/{id?}");

app.MapFallbackToFile("index.html");

app.Run();
