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

