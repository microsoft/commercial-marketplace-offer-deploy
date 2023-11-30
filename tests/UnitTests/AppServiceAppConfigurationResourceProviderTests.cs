using Microsoft.Extensions.Configuration;
using Modm.Configuration;
using Modm.Tests.Utils;

namespace Modm.Tests.UnitTests
{
    public class AppServiceAppConfigurationResourceProviderTests
    {
        private readonly string resourceGroupName;
        private readonly string subscriptionId;
        private readonly IConfiguration configuration;

        public AppServiceAppConfigurationResourceProviderTests()
        {
            resourceGroupName = Test.RandomString(20);
            subscriptionId = Guid.NewGuid().ToString();

            this.configuration = new ConfigurationBuilder()
                .AddInMemoryCollection(new Dictionary<string, string?> {
                    { AppServiceAppConfigurationResourceProvider.EnvironmentVariable_ResourceGroupName, resourceGroupName },
                    { AppServiceAppConfigurationResourceProvider.EnvironmentVariable_SubscriptionId, subscriptionId }
                }).Build();
        }

        [Fact]
        public void should_have_resource_group_match()
        {
            var provider = new AppServiceAppConfigurationResourceProvider(configuration);
            var result = provider.Get();

            Assert.Equal(resourceGroupName, result.Identifier.ResourceGroupName);
        }

        [Fact]
        public void should_have_subscription_id_match()
        {
            var provider = new AppServiceAppConfigurationResourceProvider(configuration);
            var result = provider.Get();

            Assert.Equal(subscriptionId, result.Identifier.SubscriptionId);
        }
    }
}