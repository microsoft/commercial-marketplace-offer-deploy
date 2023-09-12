using System;
using Modm.ServiceHost;

namespace Modm.ServiceHost
{
	class ControllerBuilder
	{
		readonly ControllerOptions options;
        private readonly ILogger<Startup> logger;

		IServiceProvider? serviceProvider;

        public ControllerBuilder(ILogger<Startup> logger)
		{
			this.logger = logger;
			options = new ControllerOptions();
		}

		public static ControllerBuilder Create(ILogger<Startup> logger)
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

			this.options.Logger = serviceProvider.GetRequiredService<ILogger<Controller>>();

            return new Controller(this.options);
		}
	}
}

