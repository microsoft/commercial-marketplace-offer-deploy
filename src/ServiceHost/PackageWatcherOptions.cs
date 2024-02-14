// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
namespace Modm.ServiceHost
{
	public class PackageWatcherOptions
	{
		public required string PackagePath { get; set; }
		public required string DeploymentsUrl { get; set; }
		public required string? StateFilePath { get; set; }
	}
	
}

