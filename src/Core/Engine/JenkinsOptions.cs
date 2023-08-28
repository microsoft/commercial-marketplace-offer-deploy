using System;
namespace Modm.Engine
{
	public class JenkinsOptions
	{
        public const string ConfigSectionKey = "Jenkins";

        public required string BaseUrl { get; set; }
        public required string UserName { get; set; }
        public required string ApiToken { get; set; }
	}
}

