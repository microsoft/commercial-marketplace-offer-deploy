// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using Microsoft.Extensions.Configuration;
using Microsoft.IdentityModel.Tokens;
using System.Text;

namespace Modm.Security
{
	public class JwtSettings
	{
        public class AppSettingKey
        {
            public const string SigningKey = "SigningKey";
            public const string Issuer = "Issuer";
        }

        private static readonly string Kid = "e5ccac91-4b91-4dcf-82c8-844c0fa1ddeb";

        private readonly IConfiguration configuration;

        public JwtSettings(IConfiguration configuration)
        {
            this.configuration = configuration;
        }

        public string GetIssuer()
        {
            var issuer = configuration[AppSettingKey.Issuer];

            // support local development
            if (string.IsNullOrEmpty(issuer) || issuer == "localhost")
            {
                var urls = configuration["ASPNETCORE_URLS"] ?? string.Empty;

                if (!string.IsNullOrEmpty(urls))
                {
                    return urls.Split(";").First();
                }
            }

            return issuer;
        }

        public SecurityKey GetIssuerSigningKey()
        {
            var value = configuration[AppSettingKey.SigningKey];

            if (string.IsNullOrEmpty(value))
            {
                value = Kid;
            }

            var paddedValue = value.PadRight(512, '\0');
            var key = Encoding.UTF8.GetBytes(paddedValue);

            return new SymmetricSecurityKey(key) { KeyId = Kid };
        }

        public SigningCredentials GetSigningCredentials()
        {
            return new SigningCredentials(GetIssuerSigningKey(), SecurityAlgorithms.HmacSha512Signature);
        }
    }
}

