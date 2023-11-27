using Modm.Deployments;

namespace Modm.Tests.UnitTests
{
    public class DeploymentTypeTests
	{
        [Fact]
        public void should_throw_null_argument_exception()
        {
            var exception = Record.Exception(() =>
            {
                DeploymentType nullDeploymentType = null;
            });

            Assert.True(exception is ArgumentNullException);
        }

        [Fact]
        public void should_throw_exception_when_string_is_not_valid_type()
        {
            var exception = Record.Exception(() =>
            {
                DeploymentType nullDeploymentType = "invalid";
            });

            Assert.True(exception is ArgumentOutOfRangeException);
        }

        [Fact]
        public void should_allow_implicit_conversion_of_valid_types()
        {
            foreach (var t in DeploymentType.SupportedTypes)
            {
                var exception = Record.Exception(() =>
                {
                    DeploymentType implicitlyConverted = t;
                });

                Assert.Null(exception);
            }
        }

        [Fact]
        public void jenkins_arm_deploy_script_should_use_variable_set_from_manifest_for_main_template()
        {
            var armPath = CreateDefinitionPath(DeploymentType.Arm);
            var deployScript = File.ReadAllText(Path.Combine(armPath, "deploy.sh"));

            Assert.Contains("template_file=$(cat ./manifest.json | jq -r '.mainTemplate')", deployScript);
            Assert.Contains("--template-file $template_file", deployScript);
        }

        /// <summary>
        /// Ensures that the deployment types defined match the jenkins definitions path
        /// which work through convention
        /// </summary>
        [Fact]
		public void jenkins_definitions_path_should_match_defined_types()
		{
			var arm = CreateDefinitionPath(DeploymentType.Arm);
			Assert.True(Directory.Exists(arm));

            var terraform = CreateDefinitionPath(DeploymentType.Terraform);
            Assert.True(Directory.Exists(terraform));
        }

        /// <summary>
        /// Ensures that the deployment types match the jobs created in jenkins
        /// </summary>
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

