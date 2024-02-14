// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using Microsoft.Extensions.Configuration;
using Microsoft.IdentityModel.Tokens;

namespace Modm.Security
{
	public class JwtTokenFactory
	{
        private readonly JwtSettings settings;

        public JwtTokenFactory(IConfiguration configuration)
		{
            this.settings = new JwtSettings(configuration);
		}

        public string Create(JwtTokenOptions options)
        {
            var issuer = settings.GetIssuer();
            var audience = issuer;

            var tokenDescriptor = new SecurityTokenDescriptor
            {
                Subject = new ClaimsIdentity(new[]
                {
                        new Claim("id", options.Id.ToString()),
                        new Claim(JwtRegisteredClaimNames.Sub, options.Sub),
                        new Claim(JwtRegisteredClaimNames.Jti, Guid.NewGuid().ToString())
                    }),
                Expires = options.Expires.DateTime,
                Issuer = issuer,
                Audience = audience,
                SigningCredentials = settings.GetSigningCredentials()
            };

            var tokenHandler = new JwtSecurityTokenHandler();
            var token = tokenHandler.CreateToken(tokenDescriptor);
            var jwtToken = tokenHandler.WriteToken(token);

            return jwtToken;
        }
    }
}