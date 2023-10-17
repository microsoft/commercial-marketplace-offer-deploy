using Microsoft.Extensions.DependencyInjection;
using Modm.Deployments;
using Modm.Engine;
using Modm.Engine.Pipelines;
using Modm.Extensions;
using Modm.Azure;
using Modm.Tests.Utils;
using Modm.Packaging;
using FluentValidation;
using NSubstitute;
using FluentValidation.Results;
using MediatR;

namespace Modm.Tests.UnitTests
{
    public class StartDeploymentRequestTests : AbstractTest<StartDeploymentRequestTests>
	{
		readonly StartDeploymentRequestPipeline pipeline;
        private readonly StartDeploymentRequest request;

        public StartDeploymentRequestTests() : base()
        {
            pipeline = GetPipeline();
            request = new StartDeploymentRequest
            {
                ArtifactsUri = "https://dummy-artifacts-url/installer.pkg",
                Parameters = new Dictionary<string, object>()
            };
        }

		[Fact]
		public async Task should_read_from_repository()
		{
			await pipeline.Execute(request);
            this.With<IDeploymentRepository>(async repository =>
            {
                await repository.Get(Arg.Any<CancellationToken>()).ReceivedWithAnyArgs(1);
            });
		}

        [Fact]
		public async Task should_detect_not_startable()
		{
            this.With<IDeploymentRepository>(r => r.Get().ReturnsForAnyArgs(new Deployment
            {
                IsStartable = false
            }));

            var result = await pipeline.Execute(request);

            Assert.Single(result.Errors);
            Assert.Equal("Deployment is not startable", result.Errors.First());
        }

        private StartDeploymentRequestPipeline GetPipeline()
        {
            var pipeline = Provider.GetRequiredService<IPipeline<StartDeploymentRequest, StartDeploymentResult>>();
            return (StartDeploymentRequestPipeline)pipeline;
        }

        protected override void ConfigureServices()
        {
            ConfigureMocks(m =>
            {
                m.ArtifactsDownloader();
                m.DeploymentRepository();
                m.JenkinsClient();
                m.Configuration();

                m.Create<IValidator<PackageFile>>(instance =>
                {
                    instance.Validate(Arg.Any<PackageFile>()).Returns(new ValidationResult());
                    instance.Validate(Arg.Any<IValidationContext>()).Returns(new ValidationResult());

                    Services.AddScoped(x => instance);
                });
            });

            Services.AddLogging();
            Services.AddSingleton<IMetadataService, LocalMetadataService>();
            Services.AddSingleton<IManagedIdentityService, LocalManagedIdentityService>();
            Services.AddScoped<DeploymentFile>();

            Services.AddMediatR(c => c.RegisterServicesFromAssemblyContaining<IDeploymentEngine>());
            Services.AddPipeline<IPipeline<StartDeploymentRequest, StartDeploymentResult>, StartDeploymentRequestPipeline>(c => c.AddStartDeploymentRequestPipeline());
        }
    }
}

