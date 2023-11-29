using System;
using Azure;
using Azure.Identity;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using Modm.Azure;
using NSubstitute;

namespace Modm.Tests.UnitTests
{
	public class ArmDeploymentCleanupTests
	{
		public ArmDeploymentCleanupTests()
		{
		}

        //[Fact]
        //public async Task should_delete_resources()
        //{
        //    string resourceGroupName = "rg-193-20231129094617";

        //    var azureCredentials = new DefaultAzureCredential(false);
        //    var armClient = new ArmClient(azureCredentials);
        //    var armCleanup = new AzureDeploymentCleanup(armClient);

        //    bool deleted = await armCleanup.DeleteResourcePostDeployment(resourceGroupName);
        //    Assert.True(deleted);
        //}
        [Fact]
        public async Task should_delete_resources()
        {
            // Arrange
            string resourceGroupName = "rg-193-20231129094617";

            // Create a substitute for the ArmClient
            var substituteArmClient = Substitute.For<ArmClient>(new DefaultAzureCredential(true), null);

            // Substitute the GetDefaultSubscriptionAsync method to return a fake subscription
            var fakeSubscription = Substitute.For<SubscriptionResource>();
            substituteArmClient.GetDefaultSubscriptionAsync().Returns(Task.FromResult(fakeSubscription));

            // Mock the GetResourceGroupAsync method
            var fakeResourceGroupResponse = Substitute.For<Response<ResourceGroupResource>>();
            var fakeResourceGroup = Substitute.For<ResourceGroupResource>();
            fakeResourceGroupResponse.Value.Returns(fakeResourceGroup);
            fakeSubscription.GetResourceGroupAsync(resourceGroupName).Returns(Task.FromResult(fakeResourceGroupResponse));

            // Mock the GetGenericResourcesAsync method
            fakeResourceGroup.GetGenericResourcesAsync().Returns(GetFakeAsyncEnumerable());

            // Use the substitute instead of the real ArmClient
            var armCleanup = new AzureDeploymentCleanup(substituteArmClient);

            // Act
            bool deleted = await armCleanup.DeleteResourcePostDeployment(resourceGroupName);

            // Assert
            Assert.True(deleted);
        }

        // Helper method to create a fake IAsyncEnumerable<GenericResource>
        private IAsyncEnumerable<GenericResource> GetFakeAsyncEnumerable()
        {
            // Create a substitute for IAsyncEnumerable<GenericResource>
            var fakeAsyncEnumerable = Substitute.For<IAsyncEnumerable<GenericResource>>();

            // You need to set up a way to enumerate the fakeAsyncEnumerable, depending on how your method iterates over it.
            // This is a non-trivial task and requires setting up a substitute enumerator that returns your fake resources.

            return fakeAsyncEnumerable;
        }


    }
}

