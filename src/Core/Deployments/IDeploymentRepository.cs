// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
namespace Modm.Deployments
{
	public interface IDeploymentRepository
	{
		Task<Deployment> Get(CancellationToken cancellationToken = default);
	}
}

