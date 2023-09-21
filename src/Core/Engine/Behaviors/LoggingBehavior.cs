using MediatR;
using Microsoft.Extensions.Logging;

namespace Modm.Engine.Behaviors
{
    public class LoggingBehaviour<TRequest, TResponse> : IPipelineBehavior<TRequest, TResponse>
    {
        private readonly ILogger<LoggingBehaviour<TRequest, TResponse>> logger;

        public LoggingBehaviour(ILogger<LoggingBehaviour<TRequest, TResponse>> logger)
        {
            this.logger = logger;
        }

        public async Task<TResponse> Handle(TRequest request, RequestHandlerDelegate<TResponse> next, CancellationToken cancellationToken)
        {
            logger.LogInformation("Handling {request}", typeof(TRequest).Name);

            var response = await next();

            logger.LogInformation("Handling {request}", typeof(TResponse).Name);

            return response;
        }
    }
}

