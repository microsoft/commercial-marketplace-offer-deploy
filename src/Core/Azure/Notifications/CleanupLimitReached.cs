using System;
using Azure.ResourceManager;
using MediatR;
using Modm.Deployments;

namespace Modm.Azure.Notifications
{
	public class CleanupLimitReached : IRequest
    {
        private readonly string resourceGroupName;

        public CleanupLimitReached(string resourceGroupName)
        {
            this.resourceGroupName = resourceGroupName;
        }

        public string ResourceGroupName
        {
            get { return this.resourceGroupName; }
        }
    }

    public class CleanupLimitReachedHandler : IRequestHandler<CleanupLimitReached>
    {
        private readonly IAzureResourceManager resourceManager;

        public CleanupLimitReachedHandler(IAzureResourceManager resourceManager)
        {
            this.resourceManager = resourceManager;
        }

        public async Task Handle(CleanupLimitReached request, CancellationToken cancellationToken)
        {
            var armCleanup = new AzureDeploymentCleanup(resourceManager);
            bool deleted = await armCleanup.DeleteResourcePostDeployment(request.ResourceGroupName);
        }
    }
}

