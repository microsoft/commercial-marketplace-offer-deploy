using System;
using ClientApp.Backend;
using ClientApp.Commands;
using MediatR;

namespace ClientApp.Handlers
{
    public class InitiateDeleteHandler : IRequestHandler<InitiateDelete>
    {
        private readonly DeleteService deleteService;
        private readonly ILogger<InitiateDeleteHandler> logger;

        public InitiateDeleteHandler(DeleteService deleteService, ILogger<InitiateDeleteHandler> logger)
        {
            this.deleteService = deleteService;
            this.logger = logger;
        }

        public Task Handle(InitiateDelete request, CancellationToken cancellationToken)
        {
            this.logger.LogInformation($"Handling InitiateDelete with resource group name {request.ResourceGroupName}");
            this.deleteService.Start(request.ResourceGroupName);
            return Task.CompletedTask;
        }
    }
}

