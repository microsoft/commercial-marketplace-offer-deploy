using System;
namespace ClientApp.Extensions
{
	public static class HttpRequestExtensions
	{
        /// <summary>
        /// Gets the JWT token from the request header. It is expected that the token is in the Authorization header 
        /// and that the scheme is Bearer. This MUST be the case for the proxy to work.
        /// </summary>
        /// <param name="request"></param>
        /// <returns></returns>
        public static string GetJwtToken(this HttpRequest request)
        {
            return request.Headers["Authorization"].ToString().Replace("Bearer ", "");
        }
	}
}

