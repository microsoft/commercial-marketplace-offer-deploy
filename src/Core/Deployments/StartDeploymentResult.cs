// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
namespace Modm.Deployments
{
	public record StartDeploymentResult
	{
		public Deployment Deployment { get; set; }
		public List<string> Errors { get; set; }
	}
}

