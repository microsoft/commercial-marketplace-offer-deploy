namespace Modm.ServiceHost
{
    internal class Controller
    {
        private readonly ControllerOptions options;
        private readonly ILogger<Worker> logger;

        public Controller(ControllerOptions options, ILogger<Worker> logger)
        {
            this.options = options;
            this.logger = logger;
        }

        public async Task StartAsync(CancellationToken cancellationToken = default)
        {
            logger.LogInformation("FQDN: {fqdn}", options.Fqdn);

            while (!cancellationToken.IsCancellationRequested)
            {
                logger.LogInformation("Running at: {time}", DateTimeOffset.Now);
                await Task.Delay(10000, cancellationToken);
            }
        }
    }
}