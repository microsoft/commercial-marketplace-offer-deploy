// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using ClientApp.Security;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace ClientApp.Controllers
{
    [Route("api")]
    [ApiController]
    public class AuthController : ControllerBase
	{
        private readonly AuthManager auth;
        public AuthController(AuthManager auth)
		{
            this.auth = auth;
        }

        [Route("token")]
        [AllowAnonymous]
        public async Task<IResult> Post(TokenRequest request)
        {
            var user = await auth.Get(request.Id);

            if (user != null)
            {
                return Results.Ok(user);
            }
            return Results.Unauthorized();
        }

        [Route("login")]
        [AllowAnonymous]
        public async Task<IResult> Post(LoginRequest request)
		{
            try
            {
                var result = await auth.Authenticate(request);

                if (result.IsAuthenticated)
                {
                    return Results.Ok(result.User);
                }
            }
            catch (Exception ex)
            {
                return Results.Problem(ex.Message);
            }

            return Results.Unauthorized();
        }
	}
}

