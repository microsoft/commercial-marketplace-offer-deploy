using System;
using Microsoft.WindowsAzure.ResourceStack.Common.Extensions;
using Modm.Azure.Model;

namespace Modm.Tests.UnitTests
{
	public class UserDataTests
	{
		const string Base64UserData = "eyJhcnRpZmFjdHNVcmkiOiJodHRwczovL2FtYXN0b3JhZ2Vwcm9kdXMuYmxvYi5jb3JlLndpbmRvd3MubmV0L2FwcGxpY2F0aW9uZGVmaW5pdGlvbnMvQTAwQjdfMzFFOUY5QTA5RkQyNDI5NEEwQTMwMTAxMjQ2RDk3MDBfNkE2MzM0OThBNUUzOTAyMjM4NkI3MjIzMjEwRUEwOTUwOEU0NTZCQzZFRUJGQzFDMURGN0Q4NzlCOEE5QjY4MC9jZDM3ZGIwMmE5MzM0ZWZjYmFiYzVmMDM3ZDU0ODM1MC9jb250ZW50LnppcCIsImFydGlmYWN0c0hhc2giOiIzYmU5NWY3MGYyYTIxN2FjMDI3OGNkNjJkNzJmZGYxMjczMmY5YTY5ZTkyZGI3N2YyYTcwZjVmM2U2OTk2ZTJhIiwiZnVuY3Rpb25BcHBOYW1lIjoidGVzdGZ1bmN0aW9uYXBwIiwicGFyYW1ldGVycyI6eyJsb2NhdGlvbiI6ImVhc3R1cyIsInJlc291cmNlX2dyb3VwX25hbWUiOiJ0ZXN0cmcifX0=";

        public UserDataTests()
		{
		}

		[Fact]
		public void generage_base64()
		{
            // Value of $(openssl dgst -sha256 "../Data/content.zip" | awk '{print $2}')
            string hashValue = "3be95f70f2a217ac0278cd62d72fdf12732f9a69e92db77f2a70f5f3e6996e2a";

            var parameters = new Dictionary<string, object>
            {
                { "location", "eastus" },
                { "resource_group_name", "testrg" }
            };

            var userData = new UserData
			{
				FunctionAppName = "testfunctionapp",
				ArtifactsUri = "https://amastorageprodus.blob.core.windows.net/applicationdefinitions/A00B7_31E9F9A09FD24294A0A30101246D9700_6A633498A5E39022386B7223210EA09508E456BC6EEBFC1C1DF7D879B8A9B680/cd37db02a9334efcbabc5f037d548350/content.zip",
				ArtifactsHash = hashValue,
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

			Assert.Equal(expectedUri, instance.ArtifactsUri);
		}

        [Fact]
        public void should_serialize()
        {
            var original = UserData.Deserialize(Base64UserData);

			var serialized = original.ToBase64Json();
			var deserialized = UserData.Deserialize(serialized);

            Assert.Equal(original.ArtifactsUri, deserialized.ArtifactsUri);

			Assert.Equal(original.Parameters.Count, deserialized.Parameters.Count);
			Assert.StrictEqual(original.Parameters["location"], deserialized.Parameters["location"]);
            Assert.StrictEqual(original.Parameters["resource_group_name"], deserialized.Parameters["resource_group_name"]);
        }
    }
}

