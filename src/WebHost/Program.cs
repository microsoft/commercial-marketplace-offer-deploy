using Modm;
using Modm.Artifacts;
using Modm.Engine;
using WebHost.Deployments;
using FluentValidation;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddHttpClient();

builder.Services.AddSingleton<IDeploymentEngine, JenkinsDeploymentEngine>();
builder.Services.AddSingleton<ArtifactsDownloader>();

builder.Services.AddScoped<IValidator<CreateDeploymentRequest>, CreateDeploymentRequestValidator>();

// Add services to the container.
builder.Services.AddHostedService<Worker>();
builder.Services.AddControllersWithViews();

//configuration
builder.Services.Configure<ArtifactsDownloadOptions>(builder.Configuration.GetSection(ArtifactsDownloadOptions.ConfigSectionKey));

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
