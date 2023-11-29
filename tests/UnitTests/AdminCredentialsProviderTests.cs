using Microsoft.Extensions.Configuration;
using Modm.Security;

namespace Modm.Tests.UnitTests
{
    public class AdminCredentialsProviderTests
	{
        private const string fakeBase64EncodedCredentials = "YWRtaW58a2VoaWxsaWJpY2Vwc2ltcGxlMTk0";
        private readonly IConfigurationRoot configuration;

        public AdminCredentialsProviderTests()
		{
            this.configuration = new ConfigurationBuilder()
               .AddInMemoryCollection(new Dictionary<string, string?> {
                   { AdminCredentialsProvider.AppSettingsKeyName, fakeBase64EncodedCredentials}
               }).Build();
        }

        [Fact]
        public void should_decode_appsetting()
        {
            var provider = new AdminCredentialsProvider(configuration);
            var exception = Record.Exception(() => provider.Get());
            Assert.Null(exception);
        }
    }
}