using System;
using System.IO;
using Modm.Deployments;

namespace Modm.Tests.UnitTests
{
	/// <summary>
	/// Ensures that the deployment types defined match the jenkins definitions path
	/// which work through convention
	/// </summary>
	public class DeploymentTypeTests
	{
		[Fact]
		public void jenkins_definitions_path_should_match_defined_types()
		{
			var arm = CreateDefinitionPath(DeploymentType.Arm);
			Assert.True(Directory.Exists(arm));

            var terraform = CreateDefinitionPath(DeploymentType.Terraform);
            Assert.True(Directory.Exists(terraform));
        }

        [Fact]
        public void jenkins_jobs_created_during_init_hook_should_match()
        {
            var definedJobsList = $"final jobs = ['{DeploymentType.Terraform}', '{DeploymentType.Arm}']";
            var scriptContent = GetInitHookScriptContents();

            Assert.Contains(definedJobsList, scriptContent);
        }

        private static string GetInitHookScriptContents()
        {
            var path = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "../../../../", "jenkins/init-hooks/terraform-job.groovy");
            return File.ReadAllText(path);
        }

        private static string CreateDefinitionPath(string deploymentType)
		{
            return Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "../../../../", $"jenkins/definitions/{deploymentType}");
        }
	}
}

