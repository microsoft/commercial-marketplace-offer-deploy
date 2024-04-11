using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.Extensions.Configuration;
using Microsoft.IdentityModel.Protocols.OpenIdConnect;
using Microsoft.IdentityModel.Tokens;
using Modm.Configuration;

namespace Modm.Security
{
    public class JwtBearerConfigurator
    {
        private readonly JwtSettings settings;

        public JwtBearerConfigurator(IConfiguration configuration)
        {
            this.settings = new JwtSettings(configuration);
        }

        public void Configure(JwtBearerOptions options)
        {
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
        }
    }
}

