using System;
using Microsoft.WindowsAzure.ResourceStack.Common.Extensions;
using Modm.Azure.Model;

namespace Modm.Tests.UnitTests
{
	public class UserDataTests
	{
		const string Base64UserData = "eyJkYXNoYm9hcmRVcmwiOiJ0ZXN0ZnVuY3Rpb25hcHAiLCJpbnN0YWxsZXJQYWNrYWdlIjp7InVyaSI6Imh0dHBzOi8vYW1hc3RvcmFnZXByb2R1cy5ibG9iLmNvcmUud2luZG93cy5uZXQvYXBwbGljYXRpb25kZWZpbml0aW9ucy9BMDBCN18zMUU5RjlBMDlGRDI0Mjk0QTBBMzAxMDEyNDZEOTcwMF82QTYzMzQ5OEE1RTM5MDIyMzg2QjcyMjMyMTBFQTA5NTA4RTQ1NkJDNkVFQkZDMUMxREY3RDg3OUI4QTlCNjgwL2NkMzdkYjAyYTkzMzRlZmNiYWJjNWYwMzdkNTQ4MzUwL2NvbnRlbnQuemlwIiwiaGFzaCI6IjNiZTk1ZjcwZjJhMjE3YWMwMjc4Y2Q2MmQ3MmZkZjEyNzMyZjlhNjllOTJkYjc3ZjJhNzBmNWYzZTY5OTZlMmEifSwicGFyYW1ldGVycyI6eyJsb2NhdGlvbiI6ImVhc3R1cyIsInJlc291cmNlX2dyb3VwX25hbWUiOiJ0ZXN0cmcifX0=";

        public UserDataTests()
		{
		}

        /// <summary>
        /// confirms that base64 string is properly working
        /// hash=$(openssl dgst -sha256 "../Data/content.zip" | awk '{print $2}')
        /// </summary>
        [Fact]
		public void should_generate_base64()
		{
            var parameters = new Dictionary<string, object>
            {
                { "location", "eastus" },
                { "resource_group_name", "testrg" }
            };

            var userData = new UserData
			{
				AppConfigEndpoint = "appconfigendpoint",
				InstallerPackage = new InstallerPackageInfo
				{
					Uri = "https://amastorageprodus.blob.core.windows.net/applicationdefinitions/A00B7_31E9F9A09FD24294A0A30101246D9700_6A633498A5E39022386B7223210EA09508E456BC6EEBFC1C1DF7D879B8A9B680/cd37db02a9334efcbabc5f037d548350/content.zip",
					Hash = "3be95f70f2a217ac0278cd62d72fdf12732f9a69e92db77f2a70f5f3e6996e2a"
                },
				Parameters = parameters
			};

			var base64 = userData.ToBase64Json();

			Assert.Equal(Base64UserData, base64);
		}

		[Fact]
		public void should_deserialize()
		{
			var instance = UserData.Deserialize(Base64UserData);
			var expectedUri = "https://amastorageprodus.blob.core.windows.net/applicationdefinitions/A00B7_31E9F9A09FD24294A0A30101246D9700_6A633498A5E39022386B7223210EA09508E456BC6EEBFC1C1DF7D879B8A9B680/cd37db02a9334efcbabc5f037d548350/content.zip";

			Assert.Equal(expectedUri, instance.InstallerPackage.Uri);
		}

        [Fact]
        public void should_serialize()
        {
            var original = UserData.Deserialize(Base64UserData);

			var serialized = original.ToBase64Json();
			var deserialized = UserData.Deserialize(serialized);

            Assert.Equal(original.InstallerPackage.Uri, deserialized.InstallerPackage.Uri);

			Assert.Equal(original.Parameters.Count, deserialized.Parameters.Count);
			Assert.StrictEqual(original.Parameters["location"], deserialized.Parameters["location"]);
            Assert.StrictEqual(original.Parameters["resource_group_name"], deserialized.Parameters["resource_group_name"]);
        }
    }
}

