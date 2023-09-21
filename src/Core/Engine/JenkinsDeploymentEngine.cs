using Modm.Azure;
using Modm.Deployments;
using Modm.Engine.Jenkins.Client;
using Modm.Engine.Pipelines;

namespace Modm.Engine
{
    class JenkinsDeploymentEngine : IDeploymentEngine
    {
        private readonly DeploymentFile file;
        private readonly IJenkinsClient client;
        private readonly DeploymentResourcesClient deploymentResourcesClient;
        private readonly IPipeline<StartDeploymentRequest, StartDeploymentResult> pipeline;
        private readonly IMetadataService metadataService;

        public JenkinsDeploymentEngine(DeploymentFile file, IJenkinsClient client,DeploymentResourcesClient deploymentResourcesClient,
            IPipeline<StartDeploymentRequest, StartDeploymentResult> pipeline, IMetadataService metadataService)
        {
            this.file = file;
            this.client = client;
            this.deploymentResourcesClient = deploymentResourcesClient;
            this.pipeline = pipeline;
            this.metadataService = metadataService;
        }

        public async Task<EngineInfo> GetInfo()
        {
            var info = await client.GetInfo();
            var node = await client.GetBuiltInNode();

            return new EngineInfo
            {
                EngineType = EngineType.Jenkins,
                Version = info.Version,
                IsHealthy = !node.Offline
            };
        }

        public async Task<Deployment> Get()
        {
            var deployment = await file.Read();

            // load up the resources
            var resourceGroupName = (await metadataService.GetAsync()).Compute.ResourceGroupName;
            deployment.Resources = await deploymentResourcesClient.Get(resourceGroupName);

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

