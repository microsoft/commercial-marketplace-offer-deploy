// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using Modm.Deployments;

namespace Modm.Tests.UnitTests
{
	public class TerraformParametersFileTests
	{
		public TerraformParametersFileTests()
		{
		}

		[Fact]
		public async Task values_should_serialize()
		{
			var random = Guid.NewGuid().ToString();
			var file = new TerraformParametersFile(Path.GetTempPath());
			await file.Write(new Dictionary<string, object>
			{
				{ "test", random }
			});

			var content = File.ReadAllText(Path.Combine(Path.GetTempPath(), TerraformParametersFile.FileName));

			Assert.Contains($"\"{random}\"", content);
        }
	}
}

