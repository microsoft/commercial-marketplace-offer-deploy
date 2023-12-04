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

        //[Fact]
        //public async Task DeleteResourcePostDeployment_Should_ReturnFalse_When_AnyPhaseDeletionFails()
        //{
        //    // Arrange
        //    var mockAzureResourceManager = Substitute.For<IAzureResourceManager>();
        //    var resourceGroupName = "test-rg";
        //    var phases = new string[] { "standard", "post" };
        //    var cleanup = new AzureDeploymentCleanup(mockAzureResourceManager);

        //    // Setup the mock to return a single resource for deletion
        //    mockAzureResourceManager.GetResourcesToDeleteAsync(resourceGroupName, Arg.Any<string>())
        //        .Returns(Task.FromResult(new List<GenericResource> { new GenericResource() }));

        //    // Setup the mock to simulate a deletion failure
        //    mockAzureResourceManager.TryDeleteResourceAsync(Arg.Any<GenericResource>())
        //        .Returns(Task.FromResult(false));

        //    // Act
        //    var result = await cleanup.DeleteResourcePostDeployment(resourceGroupName);

        //    // Assert
        //    Assert.False(result);
        //}

        [Fact]
        public async Task DeleteResourcePostDeployment_Should_ReturnTrue_When_AllResourcesDeleted()
        {
            // Arrange
            var mockAzureResourceManager = Substitute.For<IAzureResourceManagerClient>();
            var resourceGroupName = "test-rg";
            var cleanup = new AzureDeploymentCleanup(mockAzureResourceManager);

            // Setup the mock to return no resources for deletion
            mockAzureResourceManager.GetResourcesToDeleteAsync(resourceGroupName, Arg.Any<string>())
                .Returns(Task.FromResult(new List<GenericResource>()));

            // Act
            var result = await cleanup.DeleteResourcePostDeployment(resourceGroupName);

            // Assert
            Assert.True(result);
        }


    }
}

