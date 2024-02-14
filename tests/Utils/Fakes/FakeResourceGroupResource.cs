// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using Azure.ResourceManager.Resources;

namespace Modm.Tests.Utils.Fakes
{
	public class FakeResourceGroupResource : ResourceGroupResource
    {
		public FakeResourceGroupResource() : base()
		{
		}

		public static ResourceGroupResource New(Action<FakeResourceGroupResource>? configure = null)
		{
			var instance = new FakeResourceGroupResource();
			configure?.Invoke(instance);

			return instance;
		}
	}
}