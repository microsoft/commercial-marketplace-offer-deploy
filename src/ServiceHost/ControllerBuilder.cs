using System;
using Modm.ServiceHost;

namespace Modm.ServiceHost
{
	class ControllerBuilder
	{
		readonly ControllerOptions options;
		IServiceProvider? serviceProvider;

        public ControllerBuilder()
		{
			options = new ControllerOptions();
		}

		public static ControllerBuilder Create()
		{
			return new ControllerBuilder();
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

		public Controller Build()
		{
            if (serviceProvider == null)
            {
                throw new NullReferenceException("serviceProvider is required.");
            }

            return new Controller(this.options, serviceProvider.GetRequiredService<ILogger<Controller>>());
		}
	}
}

