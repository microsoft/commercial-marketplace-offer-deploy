using Modm.Tests.Utils;
using ClientApp.Cleanup;
using Azure.Core;
using Azure.ResourceManager;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using MediatR;
using Modm.Tests.Resources;
using Azure.ResourceManager.Resources;
using NSubstitute;
using NSubstitute.Extensions;
using Modm.Tests.Utils.Fakes;

namespace Modm.Tests.UnitTests
{
    public class DeleteProcessorTests : AbstractTest<DeleteProcessorTests>
    {
        private readonly IMediator mediator;
        private readonly string resourceGroupName;
        private readonly DeleteProcessorSpy processor;

        public DeleteProcessorTests() : base()
        {
            var client = Provider.GetRequiredService<ArmClient>();
            var logger = Provider.GetRequiredService<ILogger<DeleteProcessor>>();

            this.mediator = Provider.GetRequiredService<IMediator>();
            this.resourceGroupName = Test.RandomString(10);
            this.processor = new DeleteProcessorSpy(client, mediator, logger);
        }

        [Fact]
        public void should_have_specifically_ordered_resource_types()
        {
            var types = DeleteProcessor.ResourceTypes.Keys.ToArray();

            Assert.True(types[0].Equals(new ResourceType("Microsoft.Compute/virtualMachines")));
            Assert.True(types[1].Equals(new ResourceType("Microsoft.Network/virtualNetworks")));
            Assert.True(types[2].Equals(new ResourceType("Microsoft.Network/networkSecurityGroups")));
            Assert.True(types[3].Equals(new ResourceType("Microsoft.AppConfiguration/configurationStores")));
            Assert.True(types[4].Equals(new ResourceType("Microsoft.Storage/storageAccounts")));
            Assert.True(types[5].Equals(new ResourceType("Microsoft.Web/sites")));
        }

        [Fact]
        public void Main_Template_Has_Proper_Resources()
        {
            var expectedTypes = new string[]
            {
                "Microsoft.Storage/storageAccounts",
                "Microsoft.AppConfiguration/configurationStores",
                "Microsoft.Network/virtualNetworks",
                "Microsoft.Network/networkSecurityGroups",
                "Microsoft.Network/networkInterfaces",
                "Microsoft.Network/publicIpAddresses",
                "Microsoft.Compute/virtualMachines",
                "Microsoft.Web/serverfarms",
                "Microsoft.Web/sites"
            };

            var mainTemplate = MainTemplate.Get();
            var resources = mainTemplate.GetResourcesByCommonTag();

            Assert.Equal(expectedTypes.Length, resources.Count);
        }

        [Fact]
        public async Task should_delete_resources_from_resource_group()
        {
            await processor.DeleteResourcesAsync(resourceGroupName);

            Assert.Equal(resourceGroupName, processor.ReceivedResourceGroup);
        }

        [Fact]
        public async Task should_delete_all_resource_types()
        {
            await processor.DeleteResourcesAsync(resourceGroupName);

            await mediator.ReceivedWithAnyArgs(DeleteProcessor.ResourceTypes.Count)
                          .Send(Arg.Any<IRequest<DeleteResourceResult>>(),
                                Arg.Any<CancellationToken>());
        }

        [Fact]
        public async Task should_delete_in_order_matching_resource_types()
        {
            await processor.DeleteResourcesAsync(resourceGroupName);

            Received.InOrder(() =>
            {
                mediator.Send(Arg.Any<DeleteVirtualMachine>(), Arg.Any<CancellationToken>());
                mediator.Send(Arg.Any<DeleteVirtualNetwork>(), Arg.Any<CancellationToken>());
                mediator.Send(Arg.Any<DeleteNetworkSecurityGroup>(), Arg.Any<CancellationToken>());
                mediator.Send(Arg.Any<DeleteAppConfiguration>(), Arg.Any<CancellationToken>());
                mediator.Send(Arg.Any<DeleteStorageAccount>(), Arg.Any<CancellationToken>());
                mediator.Send(Arg.Any<DeleteAppService>(), Arg.Any<CancellationToken>());
            });
        }

        protected override void ConfigureServices()
        {
            ConfigureMocks(c =>
            {
                c.Logger<DeleteProcessor>();
                c.ArmClient();

                c.Configure(services =>
                {
                    services.AddSingleton(c.Create<IMediator>(m =>
                    {
                        m.Send(Arg.Any<IRequest<DeleteResourceResult>>(), CancellationToken.None)
                         .ReturnsForAnyArgs(new DeleteResourceResult());
                    }));
                });
            });
        }


        class DeleteProcessorSpy : DeleteProcessor
        {
            public string? ReceivedResourceGroup { get; set; }

            public DeleteProcessorSpy(ArmClient client, IMediator mediator, ILogger<DeleteProcessor> logger) : base(client, mediator, logger)
            {
            }

            protected override Task<List<GenericResource>> GetResourcesToDeleteAsync(string resourceGroupName)
            {
                ReceivedResourceGroup = resourceGroupName;
                return Task.FromResult(ResourceTypes.Select(type => FakeGenericResource.New(r =>
                {
                    r.Configure().Id.Returns(Test.AzureResourceIdentifier(resourceGroupName, type.Key));
                })).ToList());
            }

            protected override Task<ResourceGroupResource> GetResourceGroupResourceAsync(string resourceGroupName)
            {
                return Task.FromResult(FakeResourceGroupResource.New());
            }
        }
    }
}