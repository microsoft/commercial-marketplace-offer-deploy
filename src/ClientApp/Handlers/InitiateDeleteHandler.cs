using System;
using ClientApp.Backend;
using ClientApp.Commands;
using MediatR;

namespace ClientApp.Handlers
{
    public class InitiateDeleteHandler : IRequestHandler<InitiateDelete>
    {
        private readonly DeleteService deleteService;

        public InitiateDeleteHandler(DeleteService deleteService)
        {
            this.deleteService = deleteService;
        }

        public Task Handle(InitiateDelete request, CancellationToken cancellationToken)
        {
            this.deleteService.Start(request.ResourceGroupName);
            return Task.CompletedTask;
        }
    }
}

