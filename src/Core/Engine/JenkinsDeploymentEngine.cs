using System.Text;
using Microsoft.Azure.Management.ResourceManager.Fluent.Core.DAG;
using Modm.Azure;
using Modm.Deployments;
using Modm.Jenkins.Client;
using Modm.Engine.Pipelines;
using Microsoft.Extensions.Logging;

namespace Modm.Engine
{
    class JenkinsDeploymentEngine : IDeploymentEngine
    {
        private readonly DeploymentFile file;
        private readonly JenkinsClientFactory clientFactory;
        private readonly DeploymentResourcesClient deploymentResourcesClient;
        private readonly IPipeline<StartDeploymentRequest, StartDeploymentResult> pipeline;
        private readonly IPipeline<StartRedeploymentRequest, StartRedeploymentResult> redeploymentPipeline;
        private readonly IMetadataService metadataService;
        private readonly ILogger<JenkinsDeploymentEngine> logger;
        private readonly JenkinsReadinessService readinessService;

        public JenkinsDeploymentEngine(
            DeploymentFile file,
            JenkinsClientFactory clientFactory,
            DeploymentResourcesClient deploymentResourcesClient,
            IPipeline<StartDeploymentRequest,StartDeploymentResult> pipeline,
            IPipeline<StartRedeploymentRequest, StartRedeploymentResult> redeploymentPipeline,
            IMetadataService metadataService,
            JenkinsReadinessService readinessService,
            ILogger<JenkinsDeploymentEngine> logger)
        {
            this.file = file;
            this.clientFactory = clientFactory;
            this.deploymentResourcesClient = deploymentResourcesClient;
            this.pipeline = pipeline;
            this.redeploymentPipeline = redeploymentPipeline;
            this.metadataService = metadataService;
            this.readinessService = readinessService;
            this.logger = logger;
        }

        public Task<EngineInfo> GetInfo()
        {
            this.logger.LogTrace("Inside JenkinsDeploymentEngine:GetInfo()");
            return Task.FromResult(this.readinessService.GetEngineInfo());
        }

        public async Task<string> GetLogs()
        {
            using var client = await clientFactory.Create();

            var deployment = await file.ReadAsync();
            return await client.GetBuildLogs(deployment.Definition.DeploymentType, deployment.Id);
        }

        public async Task<Deployment> Get()
        {
            var deployment = await file.ReadAsync();

            // load up the resources
            var compute = (await metadataService.GetAsync()).Compute;

            deployment.SubscriptionId = compute.SubscriptionId.ToString();
            deployment.ResourceGroup = compute.ResourceGroupName;
            deployment.OfferName = compute.Offer;
            deployment.Resources = await deploymentResourcesClient.Get(compute.ResourceGroupName);

            return deployment;
        }

        /// <summary>
        /// starts a deployment
        /// </summary>
        /// <returns></returns>
        /// <remarks>
        /// see <see cref="Engine.Pipelines.StartDeploymentRequestPipeline"/> for implementation of full processing
        /// </remarks>
        public async Task<StartDeploymentResult> Start(StartDeploymentRequest request, CancellationToken cancellationToken)
        {
            this.logger.LogTrace("Inside JenkinsDeploymentEngine:Start()");
            var result = await pipeline.Execute(request, cancellationToken);
            return result;
        }

       /// <summary>
        /// starts a redeployment
        /// </summary>
        /// <returns></returns>
        /// <remarks>
        /// see <see cref="Engine.Pipelines.StartRedeploymentRequestPipeline"/> for implementation of full processing
        /// </remarks>
        public async Task<StartRedeploymentResult> Redeploy(StartRedeploymentRequest request, CancellationToken cancellationToken)
        {
            this.logger.LogTrace("Inside JenkinsDeploymentEngine:Redeploy()");
            var result = await redeploymentPipeline.Execute(request, cancellationToken);
            return result;
        }
    }
}

