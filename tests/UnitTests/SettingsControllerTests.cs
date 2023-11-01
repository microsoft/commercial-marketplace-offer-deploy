using System;
using Microsoft.Extensions.Configuration;
using Modm.ClientApp.Controllers;

namespace Modm.Tests.UnitTests
{
	public class SettingsControllerTests
	{
        private readonly IConfiguration configuration;

        public SettingsControllerTests()
		{
			this.configuration = new ConfigurationBuilder().AddInMemoryCollection(new Dictionary<string, string?>
			{
                { "url", "http://testurl" },
                { "level1", "" },
                { "level1:attr1", "attr1" },
                { "level1:attr2", "attr2" },
            }).Build();
		}

		[Fact]
		public void should_return_value_from_config()
		{
			var controller = new SettingsController(configuration);

			var result = controller.Get("url");
			const string expected = "http://testurl";

			Assert.Equal(expected, result?.ToString());
		}

        [Fact]
        public void should_support_dot_syntax_of_keys()
        {
            var controller = new SettingsController(configuration);

            var result = controller.Get("level1.attr1") as string;
   
            Assert.Equal("attr1", result);
        }

    }
}

