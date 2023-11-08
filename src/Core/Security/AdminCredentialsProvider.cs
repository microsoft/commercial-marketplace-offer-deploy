using System;
using System.Text;
using Microsoft.Extensions.Configuration;

namespace Modm.Security
{
    public class AdminCredentialsProvider
    {
        /// <summary>
        /// This value MUST be the same value set in the siteConfig from the ARM template
        /// </summary>
        /// <remarks>
        /// Example ARM template setup:
        ///     {
        ///       "name": "Credentials",
        ///       "value": "[base64(concat(parameters('_installerUsername'), '|', parameters('_installerPassword')))]"
        ///     }
        /// </remarks>
        public const string AppSettingsKeyName = "Credentials";

        private readonly IConfiguration configuration;

        public AdminCredentialsProvider(IConfiguration configuration)
        {
            this.configuration = configuration;
        }

        public AdminCredentials Get()
        {
            var base64AdminCredentials = configuration[AppSettingsKeyName];

            if (string.IsNullOrEmpty(base64AdminCredentials))
            {
                throw new InvalidCastException($"Cannot cast {AppSettingsKeyName} app setting value of null to valid AdminCredentials.");
            }

            var bytes = Convert.FromBase64String(base64AdminCredentials);
            var value = Encoding.UTF8.GetString(bytes);

            var parts = value.Split('|');

            if (parts.Length != 2)
            {
                throw new InvalidCastException($"Cannot cast {AppSettingsKeyName} app setting value to valid AdminCredentials. Invalid format.");
            }

            return new AdminCredentials { Username = parts[0], Password = parts[1] };
        }
    }
}

