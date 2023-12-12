using System.Text.Json;
using Microsoft.Extensions.DependencyInjection;
using Modm.Deployments;
using Modm.Engine;
using Modm.Engine.Pipelines;
using Modm.Extensions;
using Modm.Azure;
using Modm.Tests.Utils;
using Modm.Packaging;
using FluentValidation;
using NSubstitute;
using FluentValidation.Results;
using MediatR;
using ClientApp.Cleanup;
using Azure.ResourceManager.Resources;
using Azure.ResourceManager;
using Azure.Identity;
using Microsoft.Extensions.Logging;
using Azure.Core;

namespace Modm.Tests.UnitTests
{
	public class DeleteProcessorTests : AbstractTest<DeleteProcessorTests>
	{
		public DeleteProcessorTests()
		{
		}

        protected override void ConfigureServices()
        {
            //var mockAzureResourceManagerClient = Mock.Create<IAzureResourceManagerClient>();
            var mockMediator = Mock.Create<IMediator>();
            var mockLogger = Mock.Logger<DeleteProcessor>();

            //string standardTagName = "standard"; 
            //string postTagName = "post";

            Services.AddSingleton(mockMediator);
            Services.AddSingleton(mockLogger);

            //var credential = new DefaultAzureCredential();
            //var armClient = new ArmClient(credential);

        }

        [Fact]
        public async Task Main_Template_Has_Proper_Resources()
        {
            string[] expectedStandardTypes = new string[]
            {
                "Microsoft.Storage/storageAccounts",
                "Microsoft.AppConfiguration/configurationStores",
                "Microsoft.Network/virtualNetworks",
                "Microsoft.Network/networkSecurityGroups",
                "Microsoft.Network/networkInterfaces",
                "Microsoft.Network/publicIpAddresses",
                "Microsoft.Compute/virtualMachines"
            };

            string[] expectedPostTypes = new string[]
            {
                "Microsoft.Web/serverfarms",
                "Microsoft.Web/sites"
            };

            var testFilePath = new Uri(typeof(DeleteProcessorTests).Assembly.Location).LocalPath;
            var testDirectory = Path.GetDirectoryName(testFilePath);

            var relativePathToMainTemplate = "../../../../templates/mainTemplate.json";
            var mainTemplatePath = Path.GetFullPath(Path.Combine(testDirectory, relativePathToMainTemplate));


            var jsonContent = await File.ReadAllTextAsync(mainTemplatePath);
            Assert.NotEmpty(jsonContent);

            using var doc = JsonDocument.Parse(jsonContent);
            var root = doc.RootElement;

            var resources = root.GetProperty("resources").EnumerateArray();
            var standardResources = new List<JsonElement>();
            var postResources = new List<JsonElement>();

            FindResources(doc.RootElement, standardResources, postResources);

            Assert.NotEmpty(standardResources);
            Assert.NotEmpty(postResources);

            var standardResourceTypes = ExtractResourceTypes(standardResources);
            var postResourceTypes = ExtractResourceTypes(postResources);

            foreach (var currentStandardType in expectedStandardTypes)
            {
                Assert.Contains(currentStandardType, standardResourceTypes);
            }

            foreach (var currentPostType in expectedPostTypes)
            {
                Assert.Contains(currentPostType, postResourceTypes);
            }
        }

        private void FindResources(JsonElement element, List<JsonElement> standardResources, List<JsonElement> postResources)
        {
            if (element.ValueKind == JsonValueKind.Object)
            {
                foreach (var property in element.EnumerateObject())
                {
                    if (property.Name == "tags" && property.Value.ValueKind == JsonValueKind.String)
                    {
                        var tagReference = property.Value.GetString();
                        if (tagReference == "[variables('commonTags')]")
                        {
                            standardResources.Add(element);
                        }
                        else if (tagReference == "[variables('postTags')]")
                        {
                            postResources.Add(element);
                        }
                    }
                    else
                    {
                        FindResources(property.Value, standardResources, postResources);
                    }
                }
            }
            else if (element.ValueKind == JsonValueKind.Array)
            {
                foreach (var item in element.EnumerateArray())
                {
                    FindResources(item, standardResources, postResources);
                }
            }
        }

        private List<string> ExtractResourceTypes(List<JsonElement> resources)
        {
            var types = new List<string>();
            foreach (var resource in resources)
            {
                if (resource.TryGetProperty("type", out var typeElement) && typeElement.ValueKind == JsonValueKind.String)
                {
                    types.Add(typeElement.GetString());
                }
            }
            return types;
        }

        //[Fact]
        //public async Task DeleteInstallResourcesAsync_CallsOperationsInCorrectOrder()
        //{
        //    // Arrange
        //    var operationOrder = new List<string>();
        //    var mockMediator = Provider.GetRequiredService<IMediator>();

        //    mockMediator.Send(Arg.Do<IDeleteResourceRequest>(request =>
        //        operationOrder.Add(request.GetType().Name)))
        //        .Returns(new DeleteResourceResult { Succeeded = true });

        //    var deleteProcessor = Provider.GetRequiredService<DeleteProcessor>();

        //    // Act
        //    await deleteProcessor.DeleteInstallResourcesAsync("testResourceGroup", CancellationToken.None);

        //    // Assert
        //    var expectedOrder = new List<string>
        //    {
        //        "DeleteVirtualMachine",
        //        "DeleteVirtualNetwork",
        //        "DeleteNetworkSecurityGroup",
        //        "DeleteAppConfiguration",
        //        "DeleteStorageAccount",
        //        "DeleteAppService" // for post operation
        //    };

        //    Assert.Equal(expectedOrder, operationOrder);
        //}

    }

    //public class TestDeleteProcessor : DeleteProcessor
    //{
    //    public TestDeleteProcessor(ArmClient client, IMediator mediator, ILogger<DeleteProcessor> logger) : base(client, mediator, logger)
    //    {

    //    }

    //    protected override Task<ResourceGroupResource> GetResourceGroupResourceAsync(string resourceGroupName)
    //    {
    //        return Task.FromResult(Substitute.For<ResourceGroupResource>());
    //    }

    //    protected override Task<List<GenericResource>> GetResourcesToDeleteAsync(string resourceGroupName, string phase)
    //    {
    //        var fakeResources = new List<GenericResource>() { FakeVirtualMachineResource.New() };
    //        return Task.FromResult(fakeResources);
    //    }
    //}

    //public class FakeVirtualMachineResource : GenericResource
    //{
    //    protected FakeVirtualMachineResource() : base()
    //    {

    //    }

    //    public static FakeVirtualMachineResource New()
    //    {
    //        var fakeVirtualMachine = new FakeVirtualMachineResource();
    //        var data = fakeVirtualMachine.Data;
            
    //        fakeVirtualMachine.Data.Tags.Add("modm", "standard");
    //        fakeVirtualMachine.Data.r


    //        return fakeVirtualMachine;
    //    }

    //    public override GenericResourceData Data => FakeGenericResourceData.New();
    //}

    //public class FakeGenericResourceData : GenericResourceData
    //{

    //    protected FakeGenericResourceData() : base(new AzureLocation("eastus2"))
    //    {

    //    }

        
    //    public static FakeGenericResourceData New()
    //    {
    //        return new FakeGenericResourceData();
    //    }
    //}

    //public class FakeGenericResource : GenericResource
    //{
    //    protected FakeGenericResource() : base()
    //    {
    //        //this.
    //    }

    //    public static GenericResource New()
    //    {
    //        return new FakeGenericResource();
    //    }
    //}
}

