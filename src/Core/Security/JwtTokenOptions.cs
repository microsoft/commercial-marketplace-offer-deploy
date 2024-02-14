// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
namespace Modm.Security
{
	public class JwtTokenOptions
	{
		public Guid Id { get; set; }
		public string Sub { get; set; }
		public DateTimeOffset Expires { get; set; }
	}
}