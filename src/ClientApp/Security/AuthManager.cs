using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.IdentityModel.Tokens;
using Modm.Configuration;
using Modm.Security;

namespace ClientApp.Security
{
    public sealed class AuthManager
	{
        private readonly IMemoryCache cache;
        private readonly AdminCredentialsProvider adminCredentialsProvider;
        private readonly JwtSettings jwtSettings;

        public AuthManager(IMemoryCache cache, IConfiguration configuration, AdminCredentialsProvider adminCredentialsProvider)
		{
            this.cache = cache;
            this.adminCredentialsProvider = adminCredentialsProvider;
            this.jwtSettings = new JwtSettings(configuration);
        }

        public Task<AuthenticationResult> Authenticate(LoginRequest request)
        {
            if (MatchesAdminCredentials(request))
            {
                var user = new AuthenticatedUser();

                var issuer = jwtSettings.GetIssuer();
                var audience = issuer;

                var tokenDescriptor = new SecurityTokenDescriptor
                {
                    Subject = new ClaimsIdentity(new[]
                    {
                        new Claim("id", user.Id.ToString()),
                        new Claim(JwtRegisteredClaimNames.Sub, "Administrator"),
                        new Claim(JwtRegisteredClaimNames.Jti, Guid.NewGuid().ToString())
                    }),
                    Expires = user.Expires.DateTime,
                    Issuer = issuer,
                    Audience = audience,
                    SigningCredentials = jwtSettings.GetSigningCredentials()
                };
                var tokenHandler = new JwtSecurityTokenHandler();
                var token = tokenHandler.CreateToken(tokenDescriptor);
                var jwtToken = tokenHandler.WriteToken(token);

                user.Token = jwtToken;

                Set(user);

                return Task.FromResult(new AuthenticationResult(user));
            }

            return Task.FromResult(new AuthenticationResult(null));
        }

        public Task Set(AuthenticatedUser user)
        {
            if (cache.TryGetValue<AuthenticatedUser>(user.Id, out _))
            {
                cache.Remove(user.Id);
            }

            var expiration = user.Expires.Subtract(DateTimeOffset.UtcNow);
            var options = new MemoryCacheEntryOptions().SetSlidingExpiration(expiration);

            cache.Set(user.Id, user, options);

            return Task.CompletedTask;
        }

        public Task<AuthenticatedUser> Get(Guid id)
        {
            if (cache.TryGetValue(id, out AuthenticatedUser user))
            {
                return Task.FromResult(user);
            }
            return Task.FromResult(default(AuthenticatedUser));
        }

        private bool MatchesAdminCredentials(LoginRequest request)
        {
            var credentials = adminCredentialsProvider.Get();
            return request.Username == credentials.Username && request.Password == credentials.Password;
        }
	}
}

