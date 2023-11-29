using System;
using Azure.ResourceManager;
using MediatR;

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
        private readonly ArmClient armClient;

        public CleanupLimitReachedHandler(ArmClient armClient)
        {
            this.armClient = armClient;
        }

        public async Task Handle(CleanupLimitReached request, CancellationToken cancellationToken)
        {
            var armCleanup = new AzureDeploymentCleanup(armClient);
            bool deleted = await armCleanup.DeleteResourcePostDeployment(request.ResourceGroupName);
        }
    }
}

