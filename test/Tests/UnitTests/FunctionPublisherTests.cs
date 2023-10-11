using System;
using System.IO;
using System.IO.Compression;
using Modm.Compression;
using Xunit;
using Modm.Azure;

namespace Modm.Tests.UnitTests
{
	public class FunctionPublisherTests
	{
		public FunctionPublisherTests()
		{
		}

        [Fact]
        public async Task Azure_Function_Should_Deploy()
		{
			string functionAppName = "bobjactemplate";

            // Determine the location of the currently executing assembly.
            var currentAssemblyLocation = new Uri(System.Reflection.Assembly.GetExecutingAssembly().CodeBase).LocalPath;

            // Determine the directory from that location.
            var currentDirectory = Path.GetDirectoryName(currentAssemblyLocation);

            // Construct the relative path to the Azure Function directory.
            var relativeFunctionPath = "../../../../../src/Functions";
            var workingDirectory = Path.GetFullPath(Path.Combine(currentDirectory, relativeFunctionPath));

            if (!Directory.Exists(workingDirectory))
            {
                Assert.Fail("Working Directory Doesn't exist");
            }

            var publisher = new FunctionPublisher(workingDirectory);
			await publisher.PublishAsync(functionAppName);
        }
    }
}

