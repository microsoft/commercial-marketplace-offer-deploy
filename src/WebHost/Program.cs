using Modm;
using Modm.Artifacts;
using WebHost.Deployments;
using FluentValidation;
using Microsoft.Extensions.Azure;
using Azure.Identity;
using Modm.Extensions;
using Modm.Deployments;
using MediatR;
using Modm.Engine.Notifications;
using Modm.Engine;
using Modm.HttpClient;
using Polly.Retry;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddHttpClient();

builder.Services.AddDeploymentEngine(builder.Configuration, builder.Environment);

builder.Services.AddScoped<IValidator<CreateDeploymentRequest>, CreateDeploymentRequestValidator>();

builder.Services.AddControllersWithViews();

// azure configuration
builder.Services.AddAzureClients(clientBuilder =>
{
    clientBuilder.AddArmClient(builder.Configuration.GetSection("Azure"));
    clientBuilder.UseCredential(new DefaultAzureCredential());
});

builder.Services.AddMediatR(c =>
{
    c.RegisterServicesFromAssemblyContaining<DeploymentsController>();
});

builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowSpecificOrigin", builder =>
    {
        builder.WithOrigins(
            "https://localhost:44482",
            "https://localhost:7258",
            "https://localhost:5000",
            "http://localhost:5000");
    });
});

builder.Services.AddSingleton<AsyncRetryPolicy>(sp =>
{
    var logger = sp.GetRequiredService<ILogger<Program>>();
    return RetryPolicyProvider.GetRetryPolicy(logger);
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
app.UseCors("AllowSpecificOrigin");

app.MapControllerRoute(
    name: "default",
    pattern: "{controller}/{action=Index}/{id?}");

app.MapFallbackToFile("index.html");

app.Run();
