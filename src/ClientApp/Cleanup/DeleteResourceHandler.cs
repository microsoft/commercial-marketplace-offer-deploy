using System.Reflection;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using MediatR;
using Modm.Azure;
using Polly;

namespace ClientApp.Cleanup
{
    public abstract class DeleteResourceHandler<TRequest> : IRequestHandler<TRequest, DeleteResourceResult>
        where TRequest : IDeleteResourceRequest
    {
        protected readonly ILogger Logger;
        protected readonly ArmClient Client;

        /// <summary>
        /// The delete request being handled
        /// </summary>
        protected TRequest Request;

        public DeleteResourceHandler(ILoggerFactory loggerFactory, ArmClient client)
		{
            this.Logger = loggerFactory.CreateLogger(GetType());
            this.Client = client;
        }

        public async Task<DeleteResourceResult> Handle(TRequest request, CancellationToken cancellationToken)
        {
            Request = request;

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

        public abstract Task<DeleteResourceResult> DeleteAsync(GenericResource resource);

        protected virtual async Task<DeleteResourceResult> ExecuteAsync(TRequest request)
        {
            var resource = await GetResourceAsync(request.ResourceId);
            return await DeleteAsync(resource);
        }

        /// <summary>
        /// Gets the resource by id
        /// </summary>
        /// <param name="id"></param>
        /// <returns></returns>
        protected virtual async Task<GenericResource> GetResourceAsync(ResourceIdentifier id)
        {
            return await Client.GetGenericResource(id).GetAsync();
        }
    }
}