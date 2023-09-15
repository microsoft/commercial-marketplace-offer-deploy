using MediatR;
using Modm.ServiceHost.Notifications;

namespace Modm.ServiceHost
{
    public class ArtifactsWatcherService : BackgroundService
    {
        const int DefaultWaitDelaySeconds = 10;

        readonly ArtifactsWatcher watcher;

        ArtifactsWatcherOptions? options;
        bool controllerStarted;

        public ArtifactsWatcherService(ArtifactsWatcher watcher)
		{
            this.watcher = watcher;
        }

        protected override async Task ExecuteAsync(CancellationToken cancellationToken)
        {
            await WaitForControllerToStart(cancellationToken);

            if (options == null)
            {
                throw new InvalidOperationException("Cannot start artifacts watcher. Options are null");
            }
            await watcher.StartAsync(options);
        }

        async Task WaitForControllerToStart(CancellationToken cancellationToken)
        {
            while (!controllerStarted)
            {
                await Task.Delay(DefaultWaitDelaySeconds * 1000, cancellationToken);
            }
        }

        class ControllerStartedHandler : INotificationHandler<ControllerStarted>
        {
            private readonly ArtifactsWatcherService service;

            public ControllerStartedHandler(ArtifactsWatcherService service)
            {
                this.service = service;
            }

            public Task Handle(ControllerStarted notification, CancellationToken cancellationToken)
            {
                service.options = new ArtifactsWatcherOptions
                {
                    ArtifactsPath = notification.ArtifactsPath,
                    DeploymentsUrl = notification.DeploymentsUrl
                };
                service.controllerStarted = true;

                return Task.CompletedTask;
            }

        }
    }
}

