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
        private readonly IMetadataService metadataService;

        private readonly ILogger<JenkinsDeploymentEngine> logger;

        public JenkinsDeploymentEngine(DeploymentFile file, JenkinsClientFactory clientFactory,
            DeploymentResourcesClient deploymentResourcesClient,
            IPipeline<StartDeploymentRequest, StartDeploymentResult> pipeline, IMetadataService metadataService, ILogger<JenkinsDeploymentEngine> logger)
        {
            this.file = file;
            this.clientFactory = clientFactory;
            this.deploymentResourcesClient = deploymentResourcesClient;
            this.pipeline = pipeline;
            this.metadataService = metadataService;
            this.logger = logger;
        }

        public async Task<EngineInfo> GetInfo()
        {
            try
            {
                using var client = await clientFactory.Create();

                var info = await client.GetInfo();
                var node = await client.GetBuiltInNode();

                return new EngineInfo
                {
                    EngineType = EngineType.Jenkins,
                    Version = info.Version,
                    IsHealthy = !node.Offline
                };
            }
            catch (Exception ex)
            {
                this.logger.LogError(ex, null);
                return new EngineInfo
                {
                    EngineType = EngineType.Jenkins,
                    Version = "NA",
                    IsHealthy = false
                };
            }
        }

        public async Task<string> GetLogs()
        {
            using var client = await clientFactory.Create();

            var deployment = await file.Read();
            return await client.GetBuildLogs(deployment.Definition.DeploymentType, deployment.Id);
        }

        public async Task<Deployment> Get()
        {
            var deployment = await file.Read();
           
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
            var result = await pipeline.Execute(request, cancellationToken);
            return result;
        }
    }
}

