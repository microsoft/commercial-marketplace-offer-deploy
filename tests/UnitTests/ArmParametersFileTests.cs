// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Text.Json;
using Modm.Deployments;
using Modm.Tests.Utils;

namespace Modm.Tests.UnitTests
{
	public class ArmParametersFileTests
	{
        private readonly DeployScript deployScript;

        public ArmParametersFileTests()
		{
            this.deployScript = new DeployScript();
        }

        /// <summary>
        /// This will enforce that what the deployment script is using as the parameters file
        /// matches what we use on the backend
        /// </summary>
		[Fact]
		public void FileName_defined_in_class_should_match_deploy_script_variable_value()
		{
            Assert.Contains($"=\"{ArmParametersFile.FileName}\"", deployScript.Content);
		}

        [Fact]
        public async Task should_allow_null_parameters()
        {
            using var tempDir = Test.Directory<ArmParametersFileTests>();
            var file = new ArmParametersFile(tempDir.FullName);
            await file.Write(null);

            Assert.True(File.Exists(file.FullPath));
        }

        [Fact]
        public async Task should_write_parameters()
        {
            using var tempDir = Test.Directory<ArmParametersFileTests>();
            var file = new ArmParametersFile(tempDir.FullName);

            var parameters = new Dictionary<string, object>
            {
                {"siteName", Test.RandomString(10)},
                {"administratorLogin", Test.RandomString(12)},
                {"administratorLoginPassword", Test.RandomString(12)}
            };

            await file.Write(parameters);
            var content = JsonDocument.Parse(File.ReadAllText(file.FullPath));

            foreach (var p in parameters)
            {
                var parameter = content.RootElement.GetProperty("parameters").GetProperty(p.Key);
                Assert.Equal(p.Value, parameter.GetProperty("value").GetString());
            }
        }

        [Fact]
        public async Task should_include_required_attributes_when_written()
        {
            using var tempDir = Test.Directory<ArmParametersFileTests>();
            var file = new ArmParametersFile(tempDir.FullName);
            await file.Write(null);

            var content = JsonDocument.Parse(File.ReadAllText(file.FullPath));

            var schema = content.RootElement.GetProperty("$schema").GetString();
            Assert.Equal("https://schema.management.azure.com/schemas/2019-04-01/deploymentParameters.json#", schema);

            var contentVersion = content.RootElement.GetProperty("contentVersion").GetString();
            Assert.Equal("1.0.0.0", contentVersion);
        }

        private class DeployScript
        {
            public readonly string Content;

            public DeployScript()
            {
                var path = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "../../../../", "jenkins/definitions/arm/deploy.sh");
                this.Content = File.ReadAllText(path);
            }
        }
    }
}

