using System;
using Azure;
using Azure.Identity;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using Microsoft.Extensions.Logging;
using Modm.Azure;
using NSubstitute;

namespace Modm.Tests.UnitTests
{
    public class ArmDeploymentCleanupTests
    {
        public ArmDeploymentCleanupTests()
        {
        }

        [Fact]
        public async Task DeleteResourcePostDeployment_Should_ReturnTrue_When_AllResourcesDeleted()
        {
            // Arrange
            var mockAzureResourceManager = Substitute.For<IAzureResourceManagerClient>();
            var mockLogger = Substitute.For<ILogger<AzureDeploymentCleanup>>();
            var resourceGroupName = "test-rg";
            var cleanup = new AzureDeploymentCleanup(mockAzureResourceManager, mockLogger);

            // Setup the mock to return no resources for deletion
            mockAzureResourceManager.GetResourcesToDeleteAsync(resourceGroupName, Arg.Any<string>())
                .Returns(Task.FromResult(new List<GenericResource>()));

            // Act
            var result = await cleanup.DeleteResourcePostDeployment(resourceGroupName);

            // Assert
            Assert.True(result);
        }


        public class FakeGenericResource : GenericResource
        {
            protected FakeGenericResource(): base()
            {

            }

            public static GenericResource New()
            {
                return new FakeGenericResource();
            }
        }
    }
}

