using MediatR;

namespace Modm.ServiceHost
{
    public class ArtifactsWatcherService : BackgroundService, INotificationHandler<ControllerStarted>
    {
        const int DefaultWaitDelaySeconds = 10;

        readonly ArtifactsWatcher watcher;

        static ArtifactsWatcherOptions? options;
        static bool controllerStarted;

        public ArtifactsWatcherService(ArtifactsWatcher watcher)
		{
            this.watcher = watcher;
        }

        public Task Handle(ControllerStarted notification, CancellationToken cancellationToken)
        {
            options = new ArtifactsWatcherOptions
            {
                ArtifactsPath = notification.ArtifactsPath,
                DeploymentsUrl = notification.DeploymentsUrl
            };
            controllerStarted = true;

            return Task.CompletedTask;
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
    }
}

