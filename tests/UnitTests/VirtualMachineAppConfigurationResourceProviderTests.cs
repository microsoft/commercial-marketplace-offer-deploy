using Modm.Azure;
using Modm.Azure.Model;
using Modm.Configuration;
using Modm.Tests.Utils;
using NSubstitute;

namespace Modm.Tests.UnitTests
{
    public class VirtualMachineAppConfigurationResourceProviderTests
    {
        private readonly IMetadataService metadataService;

        public VirtualMachineAppConfigurationResourceProviderTests()
        {
            this.metadataService = Substitute.For<IMetadataService>();
        }

        [Fact]
        public void should_have_resource_group_match()
        {
            var resourceGroupName = Test.RandomString(20);
            Setup(resourceGroupName: resourceGroupName);

            var provider = new VirtualMachineAppConfigurationResourceProvider(metadataService);
            var result = provider.Get();

            Assert.Equal(resourceGroupName, result.Identifier.ResourceGroupName);
        }

        [Fact]
        public void should_have_subscription_id_match()
        {
            var subscriptionId = Guid.NewGuid();
            Setup(subscriptionId: subscriptionId);

            var provider = new VirtualMachineAppConfigurationResourceProvider(metadataService);
            var result = provider.Get();

            Assert.NotNull(result.Identifier);
            Assert.Equal(subscriptionId.ToString(), result.Identifier.SubscriptionId);
        }

        private void Setup(Guid? subscriptionId = null, string? resourceGroupName = null)
        {
            this.metadataService.GetAsync().Returns(Task.FromResult(new InstanceMetadata
            {
                Compute = new Compute
                {
                    SubscriptionId = subscriptionId.GetValueOrDefault(Guid.NewGuid()),
                    ResourceGroupName = resourceGroupName ?? Test.RandomString(20)
                },
                Network = null
            }));
        }
    }
}