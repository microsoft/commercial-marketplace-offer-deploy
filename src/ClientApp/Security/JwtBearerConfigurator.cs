using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Protocols.OpenIdConnect;
using Microsoft.IdentityModel.Tokens;

namespace ClientApp.Security
{
    public class JwtBearerConfigurator
	{
        private readonly IConfiguration configuration;

        public JwtBearerConfigurator(IConfiguration configuration)
		{
            this.configuration = configuration;
        }

		public void Configure(JwtBearerOptions options)
		{
            var settings = new JwtSettings(configuration);

            options.RequireHttpsMetadata = false;
            options.Authority = settings.GetIssuer();
            options.Configuration = new OpenIdConnectConfiguration();

            options.TokenValidationParameters = new TokenValidationParameters
            {
                ValidIssuer = settings.GetIssuer(),
                ValidAudience = settings.GetIssuer(),
                IssuerSigningKey = settings.GetIssuerSigningKey(),
                ValidateIssuer = true,
                ValidateAudience = true,
                ValidateLifetime = false,
                ValidateIssuerSigningKey = true,
            };

            options.Events = new JwtBearerEvents
            {
                OnAuthenticationFailed = context =>
                {
                    return Task.CompletedTask;
                }
            };
        }
	}
}

