//using System;
//using Azure;
//using Azure.Core;
//using Azure.ResourceManager;
//using Azure.ResourceManager.Network;
//using Azure.ResourceManager.Network.Models;
//using Azure.ResourceManager.Resources;

//namespace ClientApp.Cleanup
//{
//    public class DeleteNetworkInterface : IDeleteResourceRequest
//    {
//        public ResourceIdentifier ResourceId { get; }

//        public ResourceGroupResource ResourceGroup { get; }

//        public DeleteNetworkInterface(ResourceGroupResource resourceGroup, ResourceIdentifier identifier)
//        {
//            ResourceGroup = resourceGroup;
//            ResourceId = identifier;
//        }
//    }

//    [RetryPolicy]
//    public class DeleteNetworkInterfaceHandler : DeleteResourceHandler<DeleteNetworkInterface>
//    {
//        public DeleteNetworkInterfaceHandler(ILoggerFactory loggerFactory, ArmClient client) : base(loggerFactory, client)
//        {

//        }

//        protected override async Task<DeleteResourceResult> DeleteAsync(ResourceGroupResource resourceGroup, ResourceIdentifier id)
//        {
//            Response<NetworkInterfaceResource> response = await resourceGroup.GetNetworkInterfaceAsync(id.Name);
//            var networkInterface = response.Value;

//            var nicData = networkInterface.Data;
//            if (nicData.NetworkSecurityGroup != null)
//            {
//                var updateData = new NetworkInterfaceData()
//                {
     
//                    NetworkSecurityGroup = null,
//                    DnsSettings = networkInterface.Data.DnsSettings
//                };

//                networkInterface.Data.

//                foreach (var ipConfig in networkInterface.Data.IPConfigurations)
//                {
//                    updateData.IPConfigurations.Add(ipConfig);
//                }

//                networkInterface.u
//        }
//    }
//}

