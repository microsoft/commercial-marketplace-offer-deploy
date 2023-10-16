using System;
using Modm.Azure;
namespace Modm.Tests.UnitTests
{
	public class KuduApiTests
	{
		public KuduApiTests()
		{
		}

  //      [Fact]
		//public async Task should_deploy_zip()
		//{
  //          string functionAppName = "bobjacazinsider3";

  //          var currentAssemblyLocation = new Uri(System.Reflection.Assembly.GetExecutingAssembly().CodeBase).LocalPath;

  //          // Determine the directory from that location.
  //          var currentDirectory = Path.GetDirectoryName(currentAssemblyLocation);

  //          // Construct the relative path to the Azure Function directory.
  //          var relativeFunctionPath = "../../../Data";
  //          var dataDirectory = Path.GetFullPath(Path.Combine(currentDirectory, relativeFunctionPath));
            

  //          if (!Directory.Exists(dataDirectory))
  //          {
  //              Assert.Fail("Working Directory Doesn't exist");
  //          }

  //          var zipFilePath = Path.Combine(dataDirectory, "azurefunction.zip");
  //          var kuduApi = new KuduApi(functionAppName, new HttpClient());
  //          await kuduApi.DeployZipAsync(zipFilePath);
            
  //      }
    }
}

