using System;
using Modm.ServiceHost;

namespace Modm.ServiceHost
{
	class ControllerBuilder
	{
		readonly ControllerOptions options;
        private readonly ILogger<Worker> logger;

        public ControllerBuilder(ILogger<Worker> logger)
		{
			this.logger = logger;
			options = new ControllerOptions();
		}

		public static ControllerBuilder Create(ILogger<Worker> logger)
		{
			return new ControllerBuilder(logger);
		}

		public ControllerBuilder UseFqdn(string fqdn)
		{
			options.Fqdn = fqdn;
			return this;
		}

		public Controller Build()
		{
			return new Controller(this.options, this.logger);
		}
	}
}

