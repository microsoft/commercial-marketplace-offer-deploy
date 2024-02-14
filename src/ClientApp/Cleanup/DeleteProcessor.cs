// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System.Collections.Immutable;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using MediatR;

namespace ClientApp.Cleanup
{
    public class DeleteProcessor : IDeleteProcessor
	{
        /// <summary>
        /// The list of resource types that will be deleted in this specific order
        /// </summary>
        public static readonly Dictionary<ResourceType, Type> ResourceTypes = new(6) {
            { new ResourceType("Microsoft.Compute/virtualMachines"), typeof(DeleteVirtualMachine) },
            { new ResourceType("Microsoft.Network/virtualNetworks"), typeof(DeleteVirtualNetwork) },
            { new ResourceType("Microsoft.Network/networkSecurityGroups"), typeof(DeleteNetworkSecurityGroup) },
            { new ResourceType("Microsoft.AppConfiguration/configurationStores"), typeof(DeleteAppConfiguration) },
            { new ResourceType("Microsoft.Storage/storageAccounts"), typeof(DeleteStorageAccount) },
            { new ResourceType("Microsoft.Web/sites"), typeof(DeleteAppService) }
        };

        public const string TagName = "modm";

        private readonly ArmClient client;
        private readonly IMediator mediator;
        private readonly ILogger<DeleteProcessor> logger;
        
		public DeleteProcessor(ArmClient client, IMediator mediator, ILogger<DeleteProcessor> logger)
		{
            this.client = client;
            this.mediator = mediator;
            this.logger = logger;
		}

        public async Task DeleteResourcesAsync(string resourceGroupName, CancellationToken cancellationToken = default)
        {
            var deleteOperations = await GetDeleteOperations(resourceGroupName, cancellationToken);
            await Execute(deleteOperations);
        }

        private async Task Execute(ImmutableList<Func<Task<DeleteResourceResult>>> operations)
        {
            foreach (var operation in operations)
            {
                this.logger.LogInformation($"Executing delete command {operation}");
                var result = await operation();

                if (!result.Succeeded)
                {
                    logger.LogInformation($"Delete operation {operation} succeeded");
                }
                else
                {
                    logger.LogError($"Delete operation {operation} failed");
                }
            }
        }

        private async Task<ImmutableList<Func<Task<DeleteResourceResult>>>> GetDeleteOperations(string resourceGroupName, CancellationToken cancellationToken)
        {
            var options = new DeleteResourceOptions
            {
                ResourceGroup = await this.GetResourceGroupResourceAsync(resourceGroupName),
                Resources = await this.GetResourcesToDeleteAsync(resourceGroupName)
            };

            var deleteOperations = ResourceTypes
                .Select(item => new Func<Task<DeleteResourceResult>>(() => Delete(item.Key, item.Value, options, cancellationToken)))
                .ToImmutableList();
            return deleteOperations;
        }

        protected virtual async Task<List<GenericResource>> GetResourcesToDeleteAsync(string resourceGroupName)
        {
            var subscription = await this.client.GetDefaultSubscriptionAsync();
            var response = await subscription.GetResourceGroupAsync(resourceGroupName);
            var resourceGroup = response.Value;

            var resourcesToDelete = new List<GenericResource>();

            await foreach (var resource in resourceGroup.GetGenericResourcesAsync())
            {
                if (resource.Data.Tags != null && resource.Data.Tags.TryGetValue(TagName, out var tagValue) && tagValue == "true")
                {
                    resourcesToDelete.Add(resource);
                }
            }

            return resourcesToDelete;
        }

        protected virtual async Task<ResourceGroupResource> GetResourceGroupResourceAsync(string resourceGroupName)
        {
            var subscription = await this.client.GetDefaultSubscriptionAsync();
            var response = await subscription.GetResourceGroupAsync(resourceGroupName);
            var resourceGroup = response.Value;

            return resourceGroup;
        }

        private async Task<DeleteResourceResult> Delete(ResourceType resourceType, Type commandType, DeleteResourceOptions options, CancellationToken cancellationToken)
        {
            var command = options.CreateCommand(resourceType, commandType);

            if (command is null)
            {
                return new DeleteResourceResult { Succeeded = true };
            }

            var result = await mediator.Send(command, cancellationToken);

            if (!result.Succeeded)
            {
                logger.LogError("Delete operation {name} was not successful", commandType.Name);
            }

            return result;
        }

        public record DeleteResourceOptions
        {
            public ResourceGroupResource ResourceGroup { get; init; }
            public List<GenericResource> Resources { get; init; }

            public IDeleteResourceRequest CreateCommand(ResourceType resourceType, Type commandType)
            {
                var resource = Resources.FirstOrDefault(r => r.Id.ResourceType == resourceType);

                if (resource is null)
                {
                    return null;
                }

                return (IDeleteResourceRequest)Activator.CreateInstance(commandType, ResourceGroup, resource.Id);
            }
        }
    }
}