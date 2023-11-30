using System;
using MediatR;
using Modm.Azure;
using ClientApp.Backend;

namespace ClientApp.Notifications
{
    public class DeleteInitiated : IRequest
    {
        private readonly string resourceGroupName;

        public DeleteInitiated(string resourceGroupName)
        {
            this.resourceGroupName = resourceGroupName;
        }

        public string ResourceGroupName
        {
            get { return this.resourceGroupName; }
        }
    }

    public class DeleteInitiatedHandler : IRequestHandler<DeleteInitiated>
    {
        private readonly DeleteService deleteService;

        public DeleteInitiatedHandler(DeleteService deleteService)
        {
            this.deleteService = deleteService;
        }

        public async Task Handle(DeleteInitiated request, CancellationToken cancellationToken)
        {
            this.deleteService.Start();
        }
    }
}

