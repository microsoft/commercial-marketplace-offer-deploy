using System;
using Modm.Azure.Model;

namespace Modm.Tests
{
	public class UserDataTests
	{
		const string Base64UserData = "eyJhcnRpZmFjdHNVcmkiOiJodHRwczovL2FtYXN0b3JhZ2Vwcm9kdXMuYmxvYi5jb3JlLndpbmRvd3MubmV0L2FwcGxpY2F0aW9uZGVmaW5pdGlvbnMvQTAwQjdfMzFFOUY5QTA5RkQyNDI5NEEwQTMwMTAxMjQ2RDk3MDBfNkE2MzM0OThBNUUzOTAyMjM4NkI3MjIzMjEwRUEwOTUwOEU0NTZCQzZFRUJGQzFDMURGN0Q4NzlCOEE5QjY4MC9jZDM3ZGIwMmE5MzM0ZWZjYmFiYzVmMDM3ZDU0ODM1MC9jb250ZW50LnppcCIsInBhcmFtZXRlcnMiOnsibG9jYXRpb24iOiJlYXN0dXMiLCJyZXNvdXJjZV9ncm91cF9uYW1lIjoicmctNjQtMjAyMzA5MjExMDE5MTUifX0=";

        public UserDataTests()
		{
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
			Assert.NotStrictEqual(original.Parameters["location"], deserialized.Parameters["location"]);
            Assert.NotStrictEqual(original.Parameters["resource_group_name"], deserialized.Parameters["resource_group_name"]);
        }
    }
}

