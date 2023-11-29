using System;
using System.Text.Json;
using Microsoft.Azure.Management.ResourceManager.Fluent.Core;
using Modm.Configuration;

namespace Modm.Tests.UnitTests
{
	public class AppConfigurationResourceTests
	{
        private readonly ResourceId resourceGroupId;
        private readonly AppConfigurationResource resource;

        public AppConfigurationResourceTests()
		{
            this.resourceGroupId = ResourceId.FromString("/subscriptions/31e9f9a0-9fd2-4294-a0a3-a74fb236a071/resourceGroups/kehillitest");
            this.resource = new AppConfigurationResource(resourceGroupId);

        }

        [Fact]
		public void should_create_valid_uri()
		{
			Assert.Contains(".azconfig.io", resource.Uri.ToString());
		}

        [Fact]
        public void should_have_suffix_of_max_8_characters()
        {
			Assert.Equal(8, resource.Name.Split("-")[1].Length);
        }

        [Fact]
        public void mainTemplate_should_have_variable_that_matches()
        {
            var variable = new AppConfigStoreTemplateVariable();
            var variablesElement = variable.mainTemplate.RootElement.GetProperty("variables");

            var variableExists = variablesElement.TryGetProperty(AppConfigStoreTemplateVariable.Name, out JsonElement configStoreNameVariable);

            Assert.True(variableExists);
            Assert.Equal(AppConfigStoreTemplateVariable.Value, configStoreNameVariable.GetString());
        }

        private class AppConfigStoreTemplateVariable
        {
            /// <summary>
            /// The name of the ARM template variable that must have matching value
            /// as the what is created by the AppConfigurationResource
            /// </summary>
            public const string Name = "configStoreName";

            /// <summary>
            /// This expression matches what's replicated in AppConfigurationResource, so it
            /// MUST be this exact eval expression
            /// </summary>
            public const string Value = "[concat('modmconfig-', substring(uniqueString(resourceGroup().id), 0, 8))]";
            public readonly JsonDocument mainTemplate;

            public AppConfigStoreTemplateVariable()
            {
                var mainTemplateFilePath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "../../../../", "templates/mainTemplate.json");
                this.mainTemplate = JsonDocument.Parse(File.ReadAllText(mainTemplateFilePath));
            }
        }
    }
}

