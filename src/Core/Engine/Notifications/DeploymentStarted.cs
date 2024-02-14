// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using MediatR;

namespace Modm.Engine.Notifications
{
	public class DeploymentStarted : INotification
	{
		public int Id { get; set; }
		public string Name { get; set; }
	}
}

