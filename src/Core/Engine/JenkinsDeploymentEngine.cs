using System.Text;
using Microsoft.Azure.Management.ResourceManager.Fluent.Core.DAG;
using Modm.Azure;
using Modm.Deployments;
using Modm.Engine.Jenkins.Client;
using Modm.Engine.Pipelines;

namespace Modm.Engine
{
    class JenkinsDeploymentEngine : IDeploymentEngine
    {
        private readonly DeploymentFile file;
        private IJenkinsClient client;
        private readonly JenkinsClientFactory clientFactory;
        private readonly DeploymentResourcesClient deploymentResourcesClient;
        private readonly IPipeline<StartDeploymentRequest, StartDeploymentResult> pipeline;
        private readonly IMetadataService metadataService;

        public JenkinsDeploymentEngine(DeploymentFile file, JenkinsClientFactory clientFactory, DeploymentResourcesClient deploymentResourcesClient,
            IPipeline<StartDeploymentRequest, StartDeploymentResult> pipeline, IMetadataService metadataService)
        {
            this.file = file;
            this.clientFactory = clientFactory;
            this.client = clientFactory.Create().GetAwaiter().GetResult();
            this.deploymentResourcesClient = deploymentResourcesClient;
            this.pipeline = pipeline;
            this.metadataService = metadataService;
        }

        private IJenkinsClient JenkinsClient
        {
            get
            {
                if (this.client == null)
                {
                    this.client = this.clientFactory.Create().GetAwaiter().GetResult();
                }
                return this.client;
            }
        }

        public async Task<EngineInfo> GetInfo()
        {
            try
            {
                var info = await JenkinsClient.GetInfo();
                var node = await JenkinsClient.GetBuiltInNode();

                return new EngineInfo
                {
                    EngineType = EngineType.Jenkins,
                    Version = info.Version,
                    IsHealthy = !node.Offline
                };
            }
            catch (Exception ex)
            {
                string errorMessage = ex.Message;
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
            var deployment = await file.Read();
            return await JenkinsClient.GetBuildLogs(deployment.Definition.DeploymentType, deployment.Id);
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

