// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
using Modm.WebHost;
using Modm.Extensions;
using Modm.Azure;

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
builder.Services.AddSingleton<IAzureResourceManagerClient, AzureResourceManagerClient>();

builder.Services.AddJwtBearerAuthentication(builder.Configuration);
builder.Configuration.AddAppConfigurationSafely(builder.Environment);

var app = builder.Build();

// Configure the HTTP request pipeline.
if (!app.Environment.IsDevelopment())
{
    // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
    app.UseHsts();
}

app.UseCors("AllowLocal");

app.UseAuthentication();
app.UseAuthorization();
app.UseHttpsRedirection();

app.MapControllerRoute(
    name: "default",
    pattern: "{controller}/{action=Index}/{id?}"
).RequireAuthorization();

app.Run();
