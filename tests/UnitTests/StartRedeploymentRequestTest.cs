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
using Microsoft.WindowsAzure.ResourceStack.Common.Extensions;

namespace Modm.Tests.UnitTests
{
    public class StartRedeploymentRequestTest : AbstractTest<StartRedeploymentRequestTest>
    {
        readonly StartRedeploymentRequestPipeline pipeline;
        private readonly StartRedeploymentRequest request;

        public StartRedeploymentRequestTest() : base()
        {
            this.pipeline = GetPipeline();
            this.request = new StartRedeploymentRequest
            {
                DeploymentId = 1,
                Parameters = GetParameters()
            };

        }

        private Dictionary<string, object> GetParameters()
        {
            return new Dictionary<string, object>
            {
                { "resource_group_name", "rg-91-20230925204412" },
                { "location", "eastus" },
                { "sql_admin_password", "GPSCodeWith123" }
            };
        }
            
        private StartRedeploymentRequestPipeline GetPipeline()
        {
            var pipeline = Provider.GetRequiredService<IPipeline<StartRedeploymentRequest, StartRedeploymentResult>>();
            return (StartRedeploymentRequestPipeline)pipeline;
        }

        [Fact]
        public async Task should_submit_a_redeployment_request()
        {
            var result = await this.pipeline.Execute(request);
            Assert.NotNull(result);
        }

        protected override void ConfigureServices()
        {
            ConfigureMocks(m =>
            {
                m.PackageDownloader();
                m.DeploymentRepository();
                m.JenkinsClient();
                m.JenkinsClientFactory();
                m.DeploymentFileFactory();
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
            Services.AddSingleton<ParametersFileFactory>();
            Services.AddScoped<DeploymentFile>();
            Services.AddScoped<AuditFile>();

            Services.AddMediatR(c => c.RegisterServicesFromAssemblyContaining<IDeploymentEngine>());
            Services.AddPipeline<IPipeline<StartRedeploymentRequest, StartRedeploymentResult>, StartRedeploymentRequestPipeline>(c => c.AddStartRedeploymentRequestPipeline());
        }
    }
}

