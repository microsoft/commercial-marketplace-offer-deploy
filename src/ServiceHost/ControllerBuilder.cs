using System;
using Modm.ServiceHost;

namespace Modm.ServiceHost
{
	class ControllerBuilder
	{
		readonly ControllerOptions options;

		public ControllerBuilder()
		{
			options = new ControllerOptions();
		}

		public static ControllerBuilder Create()
		{
			return new ControllerBuilder();
		}

		public ControllerBuilder UseFqdn(string fqdn)
		{
			options.Fqdn = fqdn;
			return this;
		}

		public Controller Build()
		{
			return new Controller();
		}
	}
}

