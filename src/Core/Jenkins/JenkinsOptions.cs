// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
namespace Modm.Jenkins
{
	public class JenkinsOptions
	{
        public const string ConfigSectionKey = "Jenkins";

        public required string BaseUrl { get; set; }
        public required string UserName { get; set; }
        public required string Password { get; set; }
        public required string ApiToken { get; set; }
	}
}

