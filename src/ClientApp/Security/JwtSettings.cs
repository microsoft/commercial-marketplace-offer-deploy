using Microsoft.IdentityModel.Tokens;
using System.Text;

namespace ClientApp.Security
{
    public class JwtSettings
    {
        public const string SigningKeyAppSettingsKeyName = "SigningKey";

        /// <summary>
        /// See: https://learn.microsoft.com/en-us/azure/app-service/reference-app-settings?tabs=kudu%2Cdotnet
        /// </summary>
        public const string WebsiteNameAppSettingsKeyName = "WEBSITE_SITE_NAME";

        private static readonly string kid = "e5ccac91-4b91-4dcf-82c8-844c0fa1ddeb";

        private readonly IConfiguration configuration;

        public JwtSettings(IConfiguration configuration)
		{
            this.configuration = configuration;
        }

        public string GetIssuer()
        {
            var name = configuration[WebsiteNameAppSettingsKeyName];

            if (string.IsNullOrEmpty(name) || name == "localhost")
            {
                var urls = configuration["ASPNETCORE_URLS"] ?? string.Empty;

                if (!string.IsNullOrEmpty(urls))
                {
                    return urls.Split(";").First();
                }
            }

            return $"https://{name}.azurewebsites.net";
        }

        public SecurityKey GetIssuerSigningKey()
        {
            var value = configuration[SigningKeyAppSettingsKeyName];

            if (string.IsNullOrEmpty(value))
            {
                value = kid;
            }

            var paddedValue = value.PadRight(512, '\0');
            var key = Encoding.UTF8.GetBytes(paddedValue);

            return new SymmetricSecurityKey(key) { KeyId = kid };
        }

        public SigningCredentials GetSigningCredentials()
        {
            return new SigningCredentials(GetIssuerSigningKey(), SecurityAlgorithms.HmacSha512Signature);
        }
    }
}