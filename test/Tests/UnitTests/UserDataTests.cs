using System;
using Modm.Azure.Model;

namespace Modm.Tests.UnitTests
{
	public class UserDataTests
	{
		const string Base64UserData = "ewogICJhcnRpZmFjdHNVcmkiOiAiaHR0cHM6Ly9hbWFzdG9yYWdlcHJvZHVzLmJsb2IuY29yZS53aW5kb3dzLm5ldC9hcHBsaWNhdGlvbmRlZmluaXRpb25zL0EwMEI3XzMxRTlGOUEwOUZEMjQyOTRBMEEzMDEwMTI0NkQ5NzAwXzZBNjMzNDk4QTVFMzkwMjIzODZCNzIyMzIxMEVBMDk1MDhFNDU2QkM2RUVCRkMxQzFERjdEODc5QjhBOUI2ODAvY2QzN2RiMDJhOTMzNGVmY2JhYmM1ZjAzN2Q1NDgzNTAvY29udGVudC56aXAiLAogICJhcnRpZmFjdHNIYXNoIjogIiIsCiAgInBhcmFtZXRlcnMiOiB7CiAgICAibG9jYXRpb24iOiAiZWFzdHVzIiwKICAgICJyZXNvdXJjZV9ncm91cF9uYW1lIjogInJnLTY0LTIwMjMwOTIxMTAxOTE1IgogIH0KfQ==";

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
			Assert.StrictEqual(original.Parameters["location"], deserialized.Parameters["location"]);
            Assert.StrictEqual(original.Parameters["resource_group_name"], deserialized.Parameters["resource_group_name"]);
        }
    }
}

