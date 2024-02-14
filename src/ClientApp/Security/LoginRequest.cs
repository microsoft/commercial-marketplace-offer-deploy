// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
namespace ClientApp.Security
{
	public record LoginRequest
	{
		public string Username { get; set; }
		public string Password { get; set; }
	}
}