using System.Text;
using Microsoft.Azure.Management.ResourceManager.Fluent.Core.DAG;
using Modm.Azure;
using Modm.Deployments;
using Modm.Jenkins.Client;
using Modm.Engine.Pipelines;

namespace Modm.Engine
{
    class JenkinsDeploymentEngine : IDeploymentEngine
    {
        private readonly DeploymentFile file;
        private readonly JenkinsClientFactory clientFactory;
        private readonly DeploymentResourcesClient deploymentResourcesClient;
        private readonly IPipeline<StartDeploymentRequest, StartDeploymentResult> pipeline;
        private readonly IMetadataService metadataService;

        public JenkinsDeploymentEngine(DeploymentFile file, JenkinsClientFactory clientFactory,
            DeploymentResourcesClient deploymentResourcesClient,
            IPipeline<StartDeploymentRequest, StartDeploymentResult> pipeline, IMetadataService metadataService)
        {
            this.file = file;
            this.clientFactory = clientFactory;
            this.deploymentResourcesClient = deploymentResourcesClient;
            this.pipeline = pipeline;
            this.metadataService = metadataService;
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
            catch (Exception)
            {
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

