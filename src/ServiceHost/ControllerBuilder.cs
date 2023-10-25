using System;
using MediatR;
using Modm.Azure;
using Modm.ServiceHost;

namespace Modm.ServiceHost
{
	class ControllerBuilder
	{
		readonly ControllerOptions options;
        private readonly ILogger<ControllerService> logger;

		IServiceProvider? serviceProvider;

        public ControllerBuilder(ILogger<ControllerService> logger)
		{
			this.logger = logger;
			options = new ControllerOptions();
		}

		public static ControllerBuilder Create(ILogger<ControllerService> logger)
		{
			return new ControllerBuilder(logger);
		}

		public ControllerBuilder UsingServiceProvider(IServiceProvider serviceProvider)
		{
			this.serviceProvider = serviceProvider;
			return this;
		}

		public ControllerBuilder UseFqdn(string fqdn)
		{
			options.Fqdn = fqdn;
			return this;
		}

		public ControllerBuilder UseStateFile(string stateFile)
		{
			options.StateFilePath = stateFile;
			return this;
		}

        public ControllerBuilder UseMachineName(string machineName)
        {
            options.MachineName = machineName;
            return this;
        }

        public ControllerBuilder UseComposeFile(string composeFile)
        {
            options.ComposeFilePath = composeFile;
            return this;
        }

        public Controller Build()
		{
            if (serviceProvider == null)
            {
                throw new NullReferenceException("serviceProvider is required.");
            }

            return new Controller(this.options,
                serviceProvider.GetRequiredService<IManagedIdentityService>(),
				serviceProvider.GetRequiredService<IHostEnvironment>(),
                serviceProvider.GetRequiredService<IConfiguration>(),
                serviceProvider.GetRequiredService<IMediator>(),
                serviceProvider.GetRequiredService<ILogger<Controller>>()
            );
		}
	}
}

