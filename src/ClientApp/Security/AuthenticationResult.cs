using System;
namespace ClientApp.Security
{
    public class AuthenticationResult
    {
        public AuthenticatedUser User { get; }
        public bool IsAuthenticated { get; }

        public AuthenticationResult(AuthenticatedUser user)
        {
            User = user;
            IsAuthenticated = user != null;
        }
    }
}

