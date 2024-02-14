// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System.Reflection;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using MediatR;
using Polly;

namespace ClientApp.Cleanup
{
    public abstract class DeleteResourceHandler<TRequest> : IRequestHandler<TRequest, DeleteResourceResult>
        where TRequest : IDeleteResourceRequest
    {
        protected readonly ILogger Logger;
        protected readonly ArmClient Client;

        public DeleteResourceHandler(ILoggerFactory loggerFactory, ArmClient client)
		{
            this.Logger = loggerFactory.CreateLogger(GetType());
            this.Client = client;
        }

        public async Task<DeleteResourceResult> Handle(TRequest request, CancellationToken cancellationToken)
        {
            var retryPolicy = GetType().GetCustomAttribute<RetryPolicyAttribute>();
            if (retryPolicy == null)
            {
                return await ExecuteAsync(request);
            }

            return await Policy.Handle<Exception>()
                .WaitAndRetryAsync(retryPolicy.RetryCount, retryPolicy.GetSleepDurationProvider(), onRetry())
                .ExecuteAsync(async () => await ExecuteAsync(request));

            Action<Exception, TimeSpan, Context> onRetry()
            {
                return (ex, ts, _) =>
                {
                    Logger.LogWarning(ex, "Failed to execute handler for request {Request}, retrying after {RetryTimeSpan}s: {ExceptionMessage}",
                                    typeof(TRequest).Name, ts.TotalSeconds, ex.Message);
                };
            }
        }

        protected abstract Task<DeleteResourceResult> DeleteAsync(ResourceGroupResource resourceGroup, ResourceIdentifier id);

        protected virtual async Task<DeleteResourceResult> ExecuteAsync(TRequest request)
        {
            var subscription = await Client.GetDefaultSubscriptionAsync();
            var resourceGroup = await subscription.GetResourceGroups().GetAsync(request.ResourceId.ResourceGroupName);

            return await DeleteAsync(resourceGroup, request.ResourceId);
        }
    }
}