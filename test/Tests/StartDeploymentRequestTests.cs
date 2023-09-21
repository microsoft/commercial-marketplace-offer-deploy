using System;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Modm.Deployments;
using Modm.Engine;
using Modm.Engine.Jenkins.Client;
using Modm.Engine.Pipelines;
using NSubstitute;
using Modm.Extensions;
using Modm.Azure;
using Modm.Artifacts;
using Modm.Configuration;
using Xunit.Abstractions;
using Modm.Tests.Fakes;
using JenkinsNET.Models;
using System.Xml;
using System.Xml.Linq;

namespace Modm.Tests
{
	public class StartDeploymentRequestTests
	{
        private readonly ITestOutputHelper output;

        readonly IJenkinsClient jenkinsClient;
		readonly IConfiguration configuration;
		readonly string homePath;
		readonly StartDeploymentRequestPipeline pipeline;

		public StartDeploymentRequestTests(ITestOutputHelper output)
		{
			this.output = output;

            jenkinsClient = Substitute.For<IJenkinsClient>();
			var builds = Substitute.For<JenkinsNET.JenkinsClientBuilds>();
			builds.GetAsync<JenkinsBuildBase>(Arg.Any<string>(), Arg.Any<string>())
				.Returns(new JenkinsBuildBase(new XDocument()));

            jenkinsClient.Builds.Returns(builds);


            configuration = Substitute.For<IConfiguration>();

			homePath = Path.GetTempPath();
			output.WriteLine("Home path: " + homePath);

			configuration.GetValue<string>(EnvironmentVariable.Names.HomeDirectory).Returns(homePath);
            pipeline = GetInstance();

        }

		[Fact]
		public async Task should_read_from_disk()
		{
			var result = await pipeline.Execute(new StartDeploymentRequest
			{
				ArtifactsUri = "https://amastorageprodus.blob.core.windows.net/applicationdefinitions/A00B7_31E9F9A09FD24294A0A30101246D9700_D688E1539FCEA5BF3E9A2A0FE44B8625FE5CE8F431AD872520418A7B6116BCB8/f61384f88a23433f9b7943151d825d26/content.zip",
				Parameters = new Dictionary<string, object>()
			});

		}


		private StartDeploymentRequestPipeline GetInstance()
		{
			var services = new ServiceCollection();

			services.AddHttpClient();
            services.AddLogging();
			services.AddSingleton<IConfiguration>(configuration);

            services.AddSingleton<ArtifactsDownloader>();
            services.AddSingleton(jenkinsClient);
            services.AddSingleton<IMetadataService, LocalMetadataService>();
            services.AddSingleton<IManagedIdentityService, LocalManagedIdentityService>();
            services.AddScoped<DeploymentFile>();

            services.AddMediatR(c => c.RegisterServicesFromAssemblyContaining<IDeploymentEngine>());
            services.AddPipeline<IPipeline<StartDeploymentRequest, StartDeploymentResult>, StartDeploymentRequestPipeline>(c => c.AddStartDeploymentRequestPipeline());

            var provider = services.BuildServiceProvider();

			var instance = provider.GetRequiredService<IPipeline<StartDeploymentRequest, StartDeploymentResult>>() ?? throw new Exception();
            return (StartDeploymentRequestPipeline)instance;
        }

	}
}

